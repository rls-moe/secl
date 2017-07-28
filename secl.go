package secl

import (
	"go.rls.moe/secl/types"
	"go.rls.moe/secl/parser"
)

// ParseString turns a SECL string into a MapList
func ParseString(input string) (*types.MapList, error) {
	return parser.ParseString(input)
}

// ParseBytes turns a byteslice that contains SECL data into a MapList
func ParseBytes(input []byte) (*types.MapList, error) {
	return parser.ParseString(string(input))
}