package query

import (
	"github.com/stretchr/testify/require"
	"go.rls.moe/secl"
	"testing"
)

type UnmTest struct {
	TestString string
	TestInteger int
	TestUInteger uint
	TestFloat32 float32
	TestFloat64 float64
	TestBool bool
}

func TestQuery(t *testing.T) {
	assert := require.New(t)

	ml := secl.MustParseString(`a: (b: (c d yes) 8 9.91)`)

	test := UnmTest{}

	err := NewUnmarshalWithQuery(&test.TestString,
		NewMapKeySelect("a"),
		NewMapKeySelect("b"),
		NewListSelect(1),
	).Run(ml)

	assert.NoError(err)
	assert.Equal("d", test.TestString)

	err = NewUnmarshalWithQuery(&test.TestBool,
		NewMapKeySelect("a"),
		NewMapKeySelect("b"),
		NewListSelect(2),
	).Run(ml)

	assert.NoError(err)
	assert.True(test.TestBool)

	err = NewUnmarshalWithQuery(&test.TestInteger,
		NewMapKeySelect("a"),
		NewListSelect(0),
	).Run(ml)

	assert.NoError(err)
	assert.Equal(8, test.TestInteger)

	err = NewUnmarshalWithQuery(&test.TestUInteger,
		NewMapKeySelect("a"),
		NewListSelect(0),
	).Run(ml)

	assert.NoError(err)
	assert.Equal(uint(8), test.TestUInteger)

	err = NewUnmarshalWithQuery(&test.TestFloat32,
		NewMapKeySelect("a"),
		NewListSelect(1),
	).Run(ml)

	assert.NoError(err)
	assert.Equal(float32(9.91), test.TestFloat32)

	err = NewUnmarshalWithQuery(&test.TestFloat64,
		NewMapKeySelect("a"),
		NewListSelect(1),
	).Run(ml)

	assert.NoError(err)
	assert.Equal(9.91, test.TestFloat64)
}
