package parser

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.rls.moe/secl/types"
)

func TestParseString(t *testing.T) {

	input := "(test1:( hello test2: () test3: empty true test5: \"hallo welt\" !(nop) 88 99.1 -88 +088 -99.1 +099.1 ) off false deny no on true allow yes nil) "

	output, err := ParseString(input)
	if err != nil {
		t.Fatalf("Error on Parse: %s", err)
	}

	t.Logf("Output: %s", types.PrintValue(output))
	assert := assert.New(t)

	assert.Equal(types.PrintDebug(&types.MapList{
		Map: map[types.String]types.Value{},
		List: []types.Value{
			&types.MapList{
				Map: map[types.String]types.Value{
					types.String{Value: "test1"}: &types.MapList{
						Map: map[types.String]types.Value{
							types.String{Value: "test2"}: &types.MapList{Map: map[types.String]types.Value{}, List: []types.Value{}},
							types.String{Value: "test3"}: &types.MapList{Map: map[types.String]types.Value{}, List: []types.Value{}},
							types.String{Value: "test5"}: &types.String{Value: "hallo welt"},
						},
						List: []types.Value{
							&types.String{Value: "hello"},
							&types.Bool{Value: true},
							&types.MapList{
								Executable: true,
								Map:        map[types.String]types.Value{},
								List: []types.Value{
									&types.String{Value: "nop"},
								},
							},
							&types.Integer{Value: big.NewInt(88)},
							&types.Float{Value: big.NewFloat(99.1)},
							&types.Integer{Value: big.NewInt(-88)},
							&types.Integer{Value: big.NewInt(+88)},
							&types.Float{Value: big.NewFloat(-99.1)},
							&types.Float{Value: big.NewFloat(+99.1)},
						},
					},
				},
				List: []types.Value{
					&types.Bool{Value: false},
					&types.Bool{Value: false},
					&types.Bool{Value: false},
					&types.Bool{Value: false},
					&types.Bool{Value: true},
					&types.Bool{Value: true},
					&types.Bool{Value: true},
					&types.Bool{Value: true},
					&types.Nil{},
				},
			},
		},
	}), types.PrintDebug(output))
}
