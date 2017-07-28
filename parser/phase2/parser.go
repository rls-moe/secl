package phase2 // import "go.rls.moe/secl/parser/phase2"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/types"
)

var (
	// ErrUnbalanced is returned when a map is not properly closed
	ErrUnbalanced = errors.New("Did not find balanced map in parse")
	// ErrUnexpectedEndOfList is returned when a end of list is not expected
	ErrUnexpectedEndOfList = errors.New("Did not expect an end of list")
)

// Parser is a phase-2 parser. It turns a flat AST (or ASL) into a proper tree, folding values between two map key markers (MapBegin and MapEnd) into
// maps recursively
// It does not yet insert key-value pairs and does not perform cleanups.
type Parser struct {
	outputAST *types.MapList
	input     *phase1.AST
}

type subParser struct {
	*Parser
	depth, start, end int
}

// NewParser creates a new Phase 2 parser using the given flat AST from Phase 1
func NewParser(in *phase1.AST) *Parser {
	return &Parser{
		outputAST: &types.MapList{
			List: []types.Value{},
			Map:  map[types.String]types.Value{},
		},
		input: in,
	}
}

// Compact will transform the ASL into a AST, using the steps outlined by the documentation on the phase2.Parser struct
func (p *Parser) Compact() error {
	for x := 0; x < len(p.input.FlatNodes); x++ {
		curNode := p.input.FlatNodes[x]

		if curNode.Type() == phase1.TMapBegin {
			sp := p.subparser(1, x+1)
			var steps int
			var err error
			if steps, err = sp.compact(0); err != nil {
				return err
			}
			x += steps + 1
			p.outputAST.List = append(p.outputAST.List, sp.outputAST)
		} else if curNode.Type() == phase1.TMapEnd {
			return ErrUnexpectedEndOfList
		} else if curNode.Type() == phase1.TMapEmpty {
			p.outputAST.List = append(p.outputAST.List, &types.MapList{
				Map:  map[types.String]types.Value{},
				List: []types.Value{},
			})
		} else {
			p.outputAST.List = append(p.outputAST.List, curNode)
		}
	}

	return nil
}

// Output returns a types.MapList
//
// While this maplist may already be usable, it has not been passed through phase 3 yet so it is completely flat
// and may not be fully expanded yet. Internal types may be present.
func (p *Parser) Output() *types.MapList {
	return p.outputAST
}

func (p *Parser) subparser(depth, start int) *subParser {
	return &subParser{
		Parser: &Parser{outputAST: &types.MapList{
			List: []types.Value{},
			Map:  map[types.String]types.Value{},
		},
			input: p.input.SubAST(start, -1),
		},
		start: start,
		depth: depth,
	}
}

func (p *subParser) compact(depth int) (int, error) {
	for x := 0; x < len(p.input.FlatNodes); x++ {
		curNode := p.input.FlatNodes[x]

		if curNode.Type() == phase1.TMapBegin {
			sp := p.subparser(depth+1, x+1)
			var steps int
			var err error
			if steps, err = sp.compact(depth + 1); err != nil {
				return 0, err
			}
			x += steps + 1
			p.outputAST.List = append(p.outputAST.List, sp.outputAST)
		} else if curNode.Type() == phase1.TMapEnd {
			return x, nil
		} else if curNode.Type() == phase1.TMapEmpty {
			p.outputAST.List = append(p.outputAST.List, &types.MapList{
				Map:  map[types.String]types.Value{},
				List: []types.Value{},
			})
		} else {
			p.outputAST.List = append(p.outputAST.List, curNode)
		}
	}

	return 0, ErrUnbalanced
}
