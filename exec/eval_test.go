package exec

import (
	"github.com/stretchr/testify/require"
	"go.rls.moe/secl/parser"
	"go.rls.moe/secl/types"
	"testing"
)

func TestEval(t *testing.T) {
	assert := require.New(t)

	input := `!(nop !(nop) hel: !(nop))`

	ml, err := parser.ParseString(input)
	assert.NoError(err)

	assert.Equal(`( //MAP //LIST exec:( //MAP "hel"/STRING: exec:( //MAP //LIST nop/FUNCTION ) //LIST nop/FUNCTION exec:( //MAP //LIST nop/FUNCTION ) ) )`, types.PrintDebug(ml))

	mln, err := Eval(ml)
	assert.NoError(err)
	assert.Equal(`( //MAP //LIST nil/NIL )`, types.PrintDebug(mln))
}
