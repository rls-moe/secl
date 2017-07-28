package parser

import (
	"testing"
	"go.rls.moe/secl/types"
	assert "github.com/stretchr/testify/assert"
	"math/big"
)

func TestParseString(t *testing.T) {

	input := "(test1:( hello test2: () test3: empty true test5: \"hallo welt\" !(nop) 88 99.1 -88 +088 -99.1 +099.1 ) off false deny no on true allow yes) "

	output, err := ParseString(input)
	if err != nil {
		t.Fatalf("Error on Parse: %s", err)
	}

	t.Logf("Output: %s", types.PrintValue(output))

	assert := assert.New(t)

	assert.Equal(types.MapList{
		Map: map[types.String]types.Value{},
		List: []types.Value{
			&types.MapList{
				Map: map[types.String]types.Value{
					types.String{"test1"}: &types.MapList{
						Map: map[types.String]types.Value{
							types.String{"test2"}: &types.MapList{Map: map[types.String]types.Value{}, List:[]types.Value{}},
							types.String{"test3"}: &types.MapList{Map: map[types.String]types.Value{}, List:[]types.Value{}},
							types.String{"test5"}: &types.String{"hallo welt"},
						},
						List: []types.Value{
							&types.String{"hello"},
							&types.Bool{true},
							&types.MapList{
								Executable: true,
								Map: map[types.String]types.Value{},
								List: []types.Value{
									&types.Function{Identifier:"nop"},
								},
							},
							&types.Integer{big.NewInt(88)},
							&types.Float{big.NewFloat(99.1)},
							&types.Integer{big.NewInt(-88)},
							&types.Integer{big.NewInt(+88)},
							&types.Float{big.NewFloat(-99.1)},
							&types.Float{big.NewFloat(+99.1)},
						},
					},
				},
				List: []types.Value{
					&types.Bool{false}, &types.Bool{false},&types.Bool{false}, &types.Bool{false},
					&types.Bool{true},&types.Bool{true},&types.Bool{true},&types.Bool{true},
				},
			},
		},
	},*output)
}