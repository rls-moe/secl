package query

import (
	"go.rls.moe/secl/parser"
	"testing"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Test       string `secl:"string-test"`
	TestNum    int    `secl:"int-test"`
	TestStruct struct {
		TestStr string `secl:"a"`
		TestNum int    `secl:"b"`
	} `secl:"struct.path"`
}

func TestSimpleStructUnmarshal(t *testing.T) {
	assert := require.New(t)
	tStruct := TestConfig{}
	seclString := `
	string-test: "HI"
	int-test: 18
	struct: (
		path: (
			a: hi!
			b: 9
		)
	)
	`
	ml, err := parser.ParseString(seclString)
	assert.NoError(err)
	err = SimpleStructUnmarshal(ml, &tStruct)
	assert.NoError(err)
	assert.Equal("HI", tStruct.Test)
	assert.Equal(18, tStruct.TestNum)
	assert.Equal("hi!", tStruct.TestStruct.TestStr)
	assert.Equal(9, tStruct.TestStruct.TestNum)
}
