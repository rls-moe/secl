package secl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"go.rls.moe/secl/types"
	"math/big"
	"io/ioutil"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"strings"
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
	assert.EqualValues(types.String{Value: "world",PositionInformation: types.PositionInformation{Start:9,End:13}}, dmp.Map[types.String{Value:"hellO"}])
	assert.Len(dmp.List, 1)
	assert.EqualValues(&types.Bool{Value: false, PositionInformation: types.PositionInformation{Start:15,End:19}}, dmp.List[0])
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
	assert.EqualValues(types.String{Value: "world",PositionInformation: types.PositionInformation{Start:9,End:13}}, dmp.Map[types.String{Value:"hellO"}])
	assert.Len(dmp.List, 1)
	assert.EqualValues(&types.Bool{Value: false, PositionInformation: types.PositionInformation{Start:15,End:19}}, dmp.List[0])
}

func TestMustParse(t *testing.T) {
	assert := require.New(t)
	files, err := ioutil.ReadDir("./tests/must-parse")
	assert.NoError(err, "Must read test directory")
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".secl" {
			t.Logf("Running test %s", file.Name())
			fp := filepath.Join("./tests/must-parse", file.Name())
			data, err := ioutil.ReadFile(fp)
			assert.NoError(err, "Must read test file")
			fp2 := filepath.Join("./tests/must-parse", strings.TrimSuffix(file.Name(), ".secl") + ".expt")
			dataExpected, err := ioutil.ReadFile(fp2)
			assert.NoError(err, "Must read expected output file")

			ml, err := ParseBytes(data)
			assert.NoError(err, "Must parse without error")

			t.Logf("Output of Test %s: %s", file.Name(), types.PrintDebug(ml))

			assert.Equal(string(dataExpected), types.PrintDebug(ml), "Must match expected debug output")

			ml2, err := ParseString(types.PrintReproducableValue(ml))
			assert.NoError(err)

			assert.Equal(types.PrintValue(ml), types.PrintValue(ml2))
		}
	}
}