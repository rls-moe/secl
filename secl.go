package secl // import "go.rls.moe/secl"

import (
	"go.rls.moe/secl/types"
	"go.rls.moe/secl/parser"
)

// ParseString turns a SECL string into a MapList
func ParseString(input string) (*types.MapList, error) {
	return parser.ParseString(input)
}

// MustParseString works like ParseString but panics on error
func MustParseString(input string) (*types.MapList) {
	ml, err := ParseString(input)
	if err != nil {
		panic(err)
	}
	return ml
}

// ParseBytes turns a byteslice that contains SECL data into a MapList
func ParseBytes(input []byte) (*types.MapList, error) {
	return parser.ParseString(string(input))
}

// MustParseBytes works like ParseBytes but panics on error
func MustParseBytes(input []byte) (*types.MapList) {
	ml, err := ParseBytes(input)
	if err != nil {
		panic(err)
	}
	return ml
}