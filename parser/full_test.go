package parser

import (
	"go.rls.moe/secl/types"
	"math/big"
	"testing"
	"github.com/stretchr/testify/assert"
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
							types.String{Value: "test5"}: types.String{Value: "hallo welt", PositionInformation: types.PositionInformation{50,61}},
						},
						List: []types.Value{
							types.String{Value: "hello", PositionInformation: types.PositionInformation{9,13}},
							&types.Bool{Value: true, PositionInformation: types.PositionInformation{38,41}},
							&types.MapList{
								Executable: true,
								Map:        map[types.String]types.Value{},
								List: []types.Value{
									types.Function{Identifier: "nop"},
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
					&types.Bool{Value: false, PositionInformation: types.PositionInformation{102,104}},
					&types.Bool{Value: false, PositionInformation: types.PositionInformation{106,110}},
					&types.Bool{Value: false, PositionInformation: types.PositionInformation{112,115}},
					&types.Bool{Value: false, PositionInformation: types.PositionInformation{117,118}},
					&types.Bool{Value: true, PositionInformation: types.PositionInformation{120,121}},
					&types.Bool{Value: true, PositionInformation: types.PositionInformation{123,126}},
					&types.Bool{Value: true, PositionInformation: types.PositionInformation{128,132}},
					&types.Bool{Value: true, PositionInformation: types.PositionInformation{134,136}},
					&types.Nil{},
				},
			},
		},
	}), types.PrintDebug(output))
}
