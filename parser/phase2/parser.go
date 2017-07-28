package phase2

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/types"
)

var (
	ErrUnbalanced = errors.New("Did not find balanced map in parse")
	ErrUnexpectedEndOfList = errors.New("Did not expect an end of list")
)

type Parser struct {
	OutputAST *types.MapList
	Input     phase1.AST
}

type SubParser struct {
	*Parser
	Depth, Start, End int
}

func NewP2Parser(in phase1.AST) *Parser {
	return &Parser{
		OutputAST: &types.MapList{
			List: []types.Value{},
			Map:  map[types.String]types.Value{},
		},
		Input: in,
	}
}

func (p *Parser) Compact() error {
	for x := 0; x < len(p.Input.FlatNodes); x++ {
		curNode := p.Input.FlatNodes[x]

		if curNode.Type() == phase1.TMapBegin {
			sp := p.Subparser(1, x + 1)
			if steps, err := sp.Compact(0); err != nil {
				return err
			} else {
				x += steps+1
				p.OutputAST.List = append(p.OutputAST.List, sp.OutputAST)
			}
		} else if curNode.Type() == phase1.TMapEnd {
			return ErrUnexpectedEndOfList
		} else if curNode.Type() == phase1.TMapEmpty {
			p.OutputAST.List = append(p.OutputAST.List, &types.MapList{
				Map: map[types.String]types.Value{},
				List: []types.Value{},
			})
		} else {
			p.OutputAST.List = append(p.OutputAST.List, curNode)
		}
	}

	return nil
}

func (p *Parser) Subparser(depth, start int) *SubParser {
	return &SubParser{
		Parser: &Parser{OutputAST: &types.MapList{
			List: []types.Value{},
			Map:  map[types.String]types.Value{},
		},
			Input: &phase1.RootNode{
				p.Input.FlatNodes[start:],
			},
		},
		Start: start,
		Depth: depth,
	}
}


func (p *SubParser) Compact(depth int) (int, error) {
	for x := 0; x < len(p.Input.FlatNodes); x++ {
		curNode := p.Input.FlatNodes[x]

		if curNode.Type() == phase1.TMapBegin {
			sp := p.Subparser(depth +1, x + 1)
			if steps, err := sp.Compact(depth + 1); err != nil {
				return 0, err
			} else {
				x += steps+1
				p.OutputAST.List = append(p.OutputAST.List, sp.OutputAST)
			}
		} else if curNode.Type() == phase1.TMapEnd {
			return x, nil
		} else if curNode.Type() == phase1.TMapEmpty {
			p.OutputAST.List = append(p.OutputAST.List, &types.MapList{
				Map: map[types.String]types.Value{},
				List: []types.Value{},
			})
		} else {
			p.OutputAST.List = append(p.OutputAST.List, curNode)
		}
	}

	return 0, ErrUnbalanced
}