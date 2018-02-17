package parser // import "go.rls.moe/secl/parser"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/parser/phase2"
	"go.rls.moe/secl/parser/phase3"
	"go.rls.moe/secl/types"
)

// ParseString will parse a serialized configuration in SECL syntax into a MapList
func ParseString(input string) (*types.MapList, error) {
	defaultCtx := context.NewParserContext()

	p1 := phase1.NewParser(defaultCtx, lexer.NewTokenizer(defaultCtx, input))

	if err := p1.Run(); err != nil {
		return nil, errors.Wrap(err, "Error in Phase 1")
	}

	p2 := phase2.NewParser(p1.Output())

	if err := p2.Compact(); err != nil {
		return nil, errors.Wrap(err, "Error in Phase 2")
	}

	ast := p2.Output()
	if err := phase3.Fold(defaultCtx, ast); err != nil {
		return nil, errors.Wrap(err, "Error in Phase 3")
	}

	phase3.Clean(ast)

	return ast, nil
}
