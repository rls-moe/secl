package phase3

import (
	"testing"

	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/parser/phase2"
	"go.rls.moe/secl/types"
)

func TestPhase3Fold(t *testing.T) {
	ctx := context.NewParserContext()

	input := "(test1:( hello test2: () test3: empty true test5: hallo ) off) "

	p1 := phase1.NewParser(ctx, lexer.NewTokenizer(ctx, input))

	if err := p1.Run(); err != nil {
		t.Fatal("Error in Phase1: ", err)
	}

	p := phase2.NewParser(p1.Output())

	err := p.Compact()
	if err != nil {
		t.Fatalf("Error in step: %s", err)
	}

	out := p.Output()

	if err := Fold(ctx, out); err != nil {
		t.Fatal("Error on fold: ", err)
	}

	t.Logf("Output: %s", types.PrintValue(out))

	Clean(out)

	t.Logf("Cleand: %s", types.PrintValue(out))
}
