package phase1

import (
	"go.rls.moe/secl/types"
	"go.rls.moe/secl/lexer"
	"github.com/pkg/errors"
	"io"
	"math/rand"
)

type AST *RootNode

type Parser struct {
	FlatAST *RootNode
	Tokenizer *lexer.Tokenizer
}

func NewParser(t *lexer.Tokenizer) *Parser {
	return &Parser{
		FlatAST: &RootNode{
			FlatNodes: []types.Value{},
		},
		Tokenizer: t,
	}
}

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
	case lexer.TT_bool:
		b := &types.Bool{}
		if tok.Literal == "true" || tok.Literal == "on" || tok.Literal == "allow" || tok.Literal == "yes" {
			b.Value = true
		} else if tok.Literal == "false" || tok.Literal == "off" || tok.Literal == "deny" || tok.Literal == "no" {
			b.Value = false
		} else if tok.Literal == "maybe" {
			b.Value = rand.Float64() > 0.501
		} else {
			return errors.Errorf("Wanted a boolean value but got %q, %+v", tok.Literal, tok)
		}
		p.FlatAST.Append(b)
	case lexer.TT_empty:
		p.FlatAST.Append(&EmptyMap{})
	case lexer.TT_mapListBegin:
		p.FlatAST.Append(&MapBegin{})
	case lexer.TT_mapListEnd:
		p.FlatAST.Append(&MapEnd{})
	case lexer.TT_string:
		p.FlatAST.Append(&types.String{
			Value: tok.Literal,
		})
	case lexer.TT_singleWordString:
		p.FlatAST.Append(&types.String{
			Value: tok.Literal,
		})
	case lexer.TT_mod_execMap:
		p.FlatAST.Append(&ExecMap{})
	case lexer.TT_function:
		p.FlatAST.Append(&types.Function{
			Identifier: tok.Literal,
		})
	case lexer.TT_number:
		val, err := ConvertNumber(tok.Literal)
		if err != nil {
			return errors.Wrapf(err, "Could not convert token %q at position %d-%d",
				tok.Literal, tok.Start, tok.End)
		}
		p.FlatAST.Append(val)
	case lexer.TT_mod_mapKey:
		if err := p.FlatAST.ReplaceLast(func(in types.Value) (types.Value, error) {
			if in.Type() != types.TString {
				return nil, errors.New("Wanted string AST node")
			}
			str := in.(*types.String)
			return &MapKey{
				Value: str.Value,
			}, nil
		}); err != nil {
			return errors.Wrap(err, "Could not replace value")
		}
	case lexer.TT_eof:
		return io.EOF
	default:
		return errors.Errorf("Unknown Token %s: %+v", tok.Type, tok)
	}
	return nil
}

func (p *Parser) Output() AST {
	return AST(p.FlatAST)
}