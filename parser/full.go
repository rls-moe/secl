package parser

import (
	"go.rls.moe/secl/types"
	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/parser/phase2"
	"go.rls.moe/secl/parser/phase3"
	"github.com/pkg/errors"
)

func ParseString(input string) (*types.MapList, error) {
	p1 := phase1.NewParser(lexer.NewTokenizer(input))

	if err := p1.Run(); err != nil {
		return nil, errors.Wrap(err, "Error in Phase 1")
	}

	p2 := phase2.NewP2Parser(p1.Output())

	if err := p2.Compact(); err != nil {
		return nil, errors.Wrap(err, "Error in Phase 2")
	}

	if err := phase3.Fold(p2.OutputAST); err != nil {
		return nil, errors.Wrap(err, "Error in Phase 3")
	}

	phase3.Clean(p2.OutputAST)

	return p2.OutputAST, nil
}