package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathParse(t *testing.T) {
	assert := require.New(t)

	//secl := "d: ( () (a: ( (b:c))))"
	path := "d.[].[1].#.a.[].[0].b"
	expected := []PathSegment{
		KeySelect("d"),
		NewOnlyList(),
		ListSelect(1),
		NewOnlyMap(),
		KeySelect("a"),
		NewOnlyList(),
		ListSelect(0),
		KeySelect("b"),
	}
	parsed, err := PathToQuery(path)
	assert.NoError(err)
	assert.EqualValues(expected, parsed, "Path must parse into correct Segments")
}
