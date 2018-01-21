package secl

import (
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.rls.moe/secl/exec"
	"go.rls.moe/secl/types"
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
	assert.EqualValues(&types.String{Value: "world"}, dmp.Map[types.String{Value: "hellO"}])
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
	assert.EqualValues(&types.String{Value: "world"}, dmp.Map[types.String{Value: "hellO"}])
	assert.Len(dmp.List, 1)
	assert.EqualValues(&types.Bool{Value: false}, dmp.List[0])
}

func TestMustParse(t *testing.T) {
	assert := require.New(t)
	files, err := ioutil.ReadDir("./tests/must-parse")
	assert.NoError(err, "Must read test directory")
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".secl" {
			t.Logf("Running test %-45s: ----", file.Name())
			fp := filepath.Join("./tests/must-parse", file.Name())
			data, err := ioutil.ReadFile(fp)
			assert.NoError(err, "Must read test file")
			fp2 := filepath.Join("./tests/must-parse", strings.TrimSuffix(file.Name(), ".secl")+".expt")
			dataExpected, err := ioutil.ReadFile(fp2)
			assert.NoError(err, "Must read expected output file")

			ml, err := ParseBytes(data)
			assert.NoError(err, "Must parse without error")

			t.Logf("Output of Test %-45s: %s", file.Name(), types.PrintDebug(ml))

			assert.Equal(string(dataExpected), types.PrintDebug(ml), "Must match expected debug output")

			ml2, err := ParseString(types.PrintReproducableValue(ml))
			assert.NoError(err)

			assert.Equal(types.PrintValue(ml), types.PrintValue(ml2), "Must match reproduced value")

			if strings.HasSuffix(file.Name(), ".exec.secl") {
				ml3, err := exec.Eval(ml)
				assert.NoError(err)

				t.Logf("Output of Expanded Test %-36s: %s", file.Name(), types.PrintDebug(ml3))
				t.Logf("Parsable Output of Expanded Test  %-26s: %s", file.Name(), types.PrintReproducableValue(ml3.(*types.MapList)))

				fp3 := filepath.Join("./tests/must-parse", strings.TrimSuffix(file.Name(), ".exec.secl")+".expt")
				dataExpected, err := ioutil.ReadFile(fp3)

				assert.NoError(err, "Must read expected expanded output file")

				assert.Equal(string(dataExpected), types.PrintDebug(ml3))
			}
		}
	}
}
