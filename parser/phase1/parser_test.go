package phase1

import (
	"go.rls.moe/secl/lexer"
	"io"
	"testing"
)

func TestPhase1Parser_Step(t *testing.T) {
	input := "(test1:( hello test2: () test3: empty true test5: hallo ) off )"

	p := NewParser(lexer.NewTokenizer(input))

	for {
		err := p.Step()
		if err != nil && err != io.EOF {
			t.Fatal("Error in step: ", err)
		}
		if err == io.EOF {
			break
		}
	}

	for k := range p.FlatAST.FlatNodes {
		t.Logf("Node %d: %+#v", k, p.FlatAST.FlatNodes[k])
	}
}
