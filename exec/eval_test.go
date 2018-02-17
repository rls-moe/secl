package exec

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rls.moe/secl/parser"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/types"
)

func TestEval(t *testing.T) {
	assert := require.New(t)

	input := `!(nop !(nop) hel: !(nop))`

	ml, err := parser.ParseString(input)
	assert.NoError(err)

	assert.Equal(`( //MAP //LIST exec:( //MAP "hel"/STRING: exec:( //MAP //LIST nop/FUNCTION ) //LIST nop/FUNCTION exec:( //MAP //LIST nop/FUNCTION ) ) )`, types.PrintDebug(ml))

	ctx := context.NewParserContext()
	mln, err := Eval(ctx.ToRuntime(), ml)
	assert.NoError(err)
	assert.Equal(`( //MAP //LIST nil/NIL )`, types.PrintDebug(mln))
}
