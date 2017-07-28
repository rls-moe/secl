package phase3

import (
	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/parser/phase2"
	"go.rls.moe/secl/types"
	"testing"
)

func TestPhase3Fold(t *testing.T) {
	input := "(test1:( hello test2: () test3: empty true test5: hallo ) off) "

	p1 := phase1.NewParser(lexer.NewTokenizer(input))

	if err := p1.Run(); err != nil {
		t.Fatal("Error in Phase1: ", err)
	}

	p := phase2.NewParser(p1.Output())

	err := p.Compact()
	if err != nil {
		t.Fatalf("Error in step: %s", err)
	}

	out := p.Output()

	if err := Fold(out); err != nil {
		t.Fatal("Error on fold: ", err)
	}

	t.Logf("Output: %s", types.PrintValue(out))

	Clean(out)

	t.Logf("Cleand: %s", types.PrintValue(out))
}
