package phase1 // import "go.rls.moe/secl/parser/phase1"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/helper"
	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/types"
	"io"
	"strings"
)

// AST is a shorthand for *rootNode, an internal type, which is exported here
type AST rootNode

// SubAST returns a subset of the AST, starting at <start> and going to <end> as if the notation a.FlatNodes[<start>:<end>]
// was used. If start or end are -1, they are not inserted. Ex: start=-1 => a.FlatNodes[:<end>]
func (a *AST) SubAST(start, end int) *AST {
	var p AST
	if end != -1 && start != -1 {
		p = AST(rootNode{
			FlatNodes: a.FlatNodes[start:end],
		})
	} else if end != -1 && start == -1 {
		p = AST(rootNode{
			FlatNodes: a.FlatNodes[:end],
		})
	} else if end == -1 && start != -1 {
		p = AST(rootNode{
			FlatNodes: a.FlatNodes[start:],
		})
	} else {
		p = AST(rootNode{
			FlatNodes: a.FlatNodes[:],
		})
	}
	return &p
}

// Parser is a Phase 1 parser instance
type Parser struct {
	// FlatAST is the root node of the phase 1, a flat abstract syntax list
	FlatAST *rootNode
	// Tokenizer is the lexer to read from
	Tokenizer *lexer.Tokenizer
}

// NewParser generates a new Phase 1 parser instance from the given lexer
func NewParser(t *lexer.Tokenizer) *Parser {
	return &Parser{
		FlatAST: &rootNode{
			FlatNodes: []types.Value{},
		},
		Tokenizer: t,
	}
}

// Run will step through the parser until a io.EOF is encountered from a single step
func (p *Parser) Run() error {
	for {
		if err := p.Step(); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

// Step consumes a token from the tokenizer and turns it into a ASTNode, inserting it into the FlatAST
func (p *Parser) Step() error {
	tok := p.Tokenizer.NextToken()

	switch tok.Type {
	case lexer.TTBool:
		b := &types.Bool{}
		b.PositionInformation = types.PositionInformation{tok.Start, tok.End}
		if tok.Literal == "true" || tok.Literal == "on" || tok.Literal == "allow" || tok.Literal == "yes" {
			b.Value = true
		} else if tok.Literal == "false" || tok.Literal == "off" || tok.Literal == "deny" || tok.Literal == "no" {
			b.Value = false
		} else if tok.Literal == "maybe" {
			b.Value = helper.RndFloat() > 0.501
			b.Randomized = types.Randomized{Random: true}
		} else {
			return errors.Errorf("Wanted a boolean value but got %q, %+v", tok.Literal, tok)
		}
		if err := p.FlatAST.Append(b); err != nil {
			return err
		}
	case lexer.TTEmpty:
		if err := p.FlatAST.Append(EmptyMap{}); err != nil {
			return err
		}
	case lexer.TTNil:
		if err := p.FlatAST.Append(types.Nil{}); err != nil {
			return err
		}
	case lexer.TTMapListBegin:
		if err := p.FlatAST.Append(MapBegin{}); err != nil {
			return err
		}
	case lexer.TTMapListEnd:
		if err := p.FlatAST.Append(MapEnd{}); err != nil {
			return err
		}
	case lexer.TTString:
		replacer := strings.NewReplacer(
			"\\n", "\n",
			"\\t", "\t",
			"\\\"", "\"",
			"\\\\", "\\",
		)
		tok.Literal = replacer.Replace(tok.Literal)
		if err := p.FlatAST.Append(types.String{
			Value: tok.Literal,
			PositionInformation: types.PositionInformation{
				Start: tok.Start, End: tok.End,
			},
		}); err != nil {
			return err
		}
	case lexer.TTSingleWordString:
		if err := p.FlatAST.Append(types.String{
			Value: tok.Literal,
			PositionInformation: types.PositionInformation{
				Start: tok.Start, End: tok.End,
			},
		}); err != nil {
			return nil
		}
	case lexer.TTModExecMap:
		p.FlatAST.Append(ExecMap{})
	case lexer.TTFunction:
		if err := p.FlatAST.Append(types.Function{
			Identifier: tok.Literal,
		}); err != nil {
			return err
		}
	case lexer.TTNumber:
		val, err := ConvertNumber(tok.Literal)
		if err != nil {
			return errors.Wrapf(err, "Could not convert token %q at position %d-%d",
				tok.Literal, tok.Start, tok.End)
		}
		if err := p.FlatAST.Append(val); err != nil {
			return err
		}
	case lexer.TTModMapKey:
		if err := p.FlatAST.ReplaceLast(func(in types.Value) (types.Value, error) {
			if in.Type() != types.TString {
				return nil, errors.New("Wanted string AST node")
			}
			str := in.(types.String)
			return &MapKey{
				Value: str,
			}, nil
		}); err != nil {
			return errors.Wrap(err, "Could not replace value")
		}
	case lexer.TTRandstr:
		var length = 42
		if strings.HasSuffix(tok.Literal, "32") {
			length = 32
		} else if strings.HasSuffix(tok.Literal, "64") {
			length = 64
		} else if strings.HasSuffix(tok.Literal, "128") {
			length = 64
		} else if strings.HasSuffix(tok.Literal, "256") {
			length = 64
		}
		if err := p.FlatAST.Append(types.String{
			Value:               helper.RndStr(length),
			PositionInformation: types.PositionInformation{tok.Start, tok.End},
			Randomized:          types.Randomized{Random: true},
		}); err != nil {
			return err
		}
	case lexer.TTEOF:
		return io.EOF
	case lexer.TTModTrim:
		p.FlatAST.ModNext = func(value types.Value) (types.Value, error) {
			if value.Type() != types.TString {
				return nil, errors.Errorf("TrimString modification was not followed by a string but by %s", value.Type())
			}
			return trimString(value.(types.String)), nil
		}
	default:
		return errors.Errorf("Unknown Token %s: %+v", tok.Type, tok)
	}
	return nil
}

// Output returns the Phase 1 AST
//
// In this phase, the AST is purely flat and not an actual tree, it's more a ASL, Abstract Syntax List
//
// However, this step is mainly concerned with parsing the tokens into the correct types or preparing them
// for the next phases
func (p *Parser) Output() *AST {
	q := (*AST)(p.FlatAST)
	return q
}

func trimString(s types.String) types.String {
	s.Value = strings.TrimSpace(s.Value)
	return s
}
