package secl // import "go.rls.moe/secl"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/exec"
	"go.rls.moe/secl/parser"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/query"
	"go.rls.moe/secl/types"
)

// ParseString turns a SECL string into a MapList
func ParseString(input string) (*types.MapList, error) {
	return parser.ParseString(input)
}

// MustParseString works like ParseString but panics on error
func MustParseString(input string) *types.MapList {
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
func MustParseBytes(input []byte) *types.MapList {
	ml, err := ParseBytes(input)
	if err != nil {
		panic(err)
	}
	return ml
}

// LoadConfig will read a given string and expand this. When using loadf and loadd, you can use a hardcoded
// config and manually loading from disk is not necessary
// This function will also strip the root map from the config so you can address everything directly
func LoadConfig(config string) (*types.MapList, error) {
	cfg, err := ParseString(config)
	if err != nil {
		return nil, err
	}
	ctx := context.NewParserContext()
	ncfg, err := exec.Eval(ctx.ToRuntime(), cfg)
	if err != nil {
		return nil, err
	}
	ncfg, err = query.NewQuery(query.ListSelect(0)).Run(ncfg.(*types.MapList))
	if err != nil {
		return nil, err
	}
	cfg, ok := ncfg.(*types.MapList)
	if !ok {
		return nil, errors.New("Could not cast to MapList")
	}
	return cfg, nil
}
