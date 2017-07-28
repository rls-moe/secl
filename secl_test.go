package secl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"go.rls.moe/secl/types"
	"math/big"
)

func TestParseBytes(t *testing.T) {
	assert := assert.New(t)

	mapList, err := ParseBytes([]byte("( hellO: world false ) 99"))

	assert.NoError(err, "No error in parse")

	assert.Len(mapList.List, 2)
	assert.Len(mapList.Map, 0)
	assert.EqualValues(mapList.List[1], &types.Integer{Value: big.NewInt(99)}, "Assert List Entry Equality")

	dmp := mapList.List[0].(*types.MapList)

	assert.Len(dmp.Map, 1)
	assert.EqualValues(&types.String{Value: "world"}, dmp.Map[types.String{Value:"hellO"}])
	assert.Len(dmp.List, 1)
	assert.EqualValues(&types.Bool{Value: false}, dmp.List[0])
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)

	mapList, err := ParseString("( hellO: world false ) 99")

	assert.NoError(err, "No error in parse")

	assert.Len(mapList.List, 2)
	assert.Len(mapList.Map, 0)
	assert.EqualValues(mapList.List[1], &types.Integer{Value: big.NewInt(99)}, "Assert List Entry Equality")

	dmp := mapList.List[0].(*types.MapList)

	assert.Len(dmp.Map, 1)
	assert.EqualValues(&types.String{Value: "world"}, dmp.Map[types.String{Value:"hellO"}])
	assert.Len(dmp.List, 1)
	assert.EqualValues(&types.Bool{Value: false}, dmp.List[0])
}
