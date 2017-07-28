package phase2

import (
	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser/phase1"
	"testing"
	"go.rls.moe/secl/types"
)

func TestPhase2Parser_Step(t *testing.T) {
	input := "(test1:( hello test2: () test3: empty true test5: hallo ) off) "

	p1 := phase1.NewParser(lexer.NewTokenizer(input))

	if err := p1.Run(); err != nil {
		t.Fatalf("Error in Phase1: ", err)
	}

	p := NewP2Parser(p1.Output())

	err := p.Compact()
	if err != nil {
		t.Fatalf("Error in step: %s", err)
	}

	for k := range p.OutputAST.List {
		t.Logf("Node %d: %+#v", k, p.OutputAST.List[k])
	}

	t.Logf("Output %s", types.PrintValue(p.OutputAST))

	t.Logf("Found %d items", len(p.OutputAST.List))
}
