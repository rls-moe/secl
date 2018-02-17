package exec // import "go.rls.moe/secl/exec"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/types"
)

// Errors returned by the EvalMapList function when running an executable map
var (
	ErrListLengthZero  = errors.New("Length of List was zero, cannot be a function call")
	ErrNotAFunction    = errors.New("First list element was not a function")
	ErrUnknownFunction = errors.New("First list element was a function but it's unknown")
	ErrNotExecutable   = errors.New("MapList is not marked executable")
)

func init() {
	fn := context.MustRegisterFunction
	fn("nop", func(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
		return types.Nil{}, nil
	})
	fn("decb64", decb64)
	fn("loadv", loadv)
	fn("loadb", loadb)
	fn("loadf", loadf)
	fn("loadd", loadd)
	fn("merge", merge)
}

// EvalMapList executes a MapList which has been marked executable with the correct function
// The first element of the list must be a Function type with a valid identifier
func EvalMapList(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
	if !list.Executable {
		return nil, ErrNotExecutable
	}
	if len(list.List) == 0 {
		return nil, ErrListLengthZero
	}
	if list.List[0].Type() != types.TFunction {
		return nil, errors.Errorf("First element was not a function: %s", types.PrintValue(list.List[0]))
	}
	fnc := list.List[0].(types.Function)
	fncCall, ok := ctx.Functions[fnc.Identifier]
	if !ok {
		return nil, ErrUnknownFunction
	}

	return fncCall(ctx, list)
}

/*
func RegisterFunction(ctx context.Parser, keyword string, fnc SECLFunc) error {
	if err := lexer.RegisterFunctionKeyword(keyword); err != nil {
		return err
	}
	ctx.Functions[keyword] = fnc
	return nil
}
*/
