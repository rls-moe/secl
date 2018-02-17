package phase2

import (
	"testing"

	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/types"
)

func TestPhase2Parser_Step(t *testing.T) {
	ctx := context.NewParserContext()

	input := "(test1:( hello test2: () test3: empty true test5: hallo ) off) "

	p1 := phase1.NewParser(ctx, lexer.NewTokenizer(ctx, input))

	if err := p1.Run(); err != nil {
		t.Fatal("Error in Phase1: ", err)
	}

	p := NewParser(p1.Output())

	err := p.Compact()
	if err != nil {
		t.Fatalf("Error in step: %s", err)
	}

	for k := range p.outputAST.List {
		t.Logf("Node %d: %+#v", k, p.outputAST.List[k])
	}

	t.Logf("Output %s", types.PrintValue(p.outputAST))

	t.Logf("Found %d items", len(p.outputAST.List))
}
