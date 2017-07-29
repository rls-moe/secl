package secl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"go.rls.moe/secl/types"
	"math/big"
	"io/ioutil"
	"github.com/stretchr/testify/require"
	"path/filepath"
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
	assert.EqualValues(types.String{Value: "world"}, dmp.Map[types.String{Value:"hellO"}])
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
	assert.EqualValues(types.String{Value: "world"}, dmp.Map[types.String{Value:"hellO"}])
	assert.Len(dmp.List, 1)
	assert.EqualValues(&types.Bool{Value: false}, dmp.List[0])
}

func TestMustParse(t *testing.T) {
	assert := require.New(t)
	files, err := ioutil.ReadDir("./tests/must-parse")
	assert.NoError(err, "Must read test directory")
	for _, file := range files {
		t.Logf("Running test %s", file.Name())
		fp := filepath.Join("./tests/must-parse", file.Name())
		data, err := ioutil.ReadFile(fp)
		assert.NoError(err, "Must read test file")
		ml, err := ParseBytes(data)
		assert.NoError(err, "Must parse without error")
		t.Logf("Output of test %s:\n%s", file.Name(), types.PrintValue(ml))
	}
}