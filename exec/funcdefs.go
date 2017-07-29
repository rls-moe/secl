package exec // import "go.rls.moe/secl/exe"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
)

// SECLFunc is a generic function to be executed in a SECL file
// It receives the maplist that contained the function
// When it returns no error, it must return a non-nil types.Value entity that
// replaces/expands the function position
type SECLFunc func(list *types.MapList) (types.Value, error)

// Errors returned by the EvalMapList function when running an executable map
var (
	ErrListLengthZero  = errors.New("Length of List was zero, cannot be a function call")
	ErrNotAFunction    = errors.New("First list element was not a function")
	ErrUnknownFunction = errors.New("First list element was a function but it's unknown")
	ErrNotExecutable   = errors.New("MapList is not marked executable")
)

var functions = map[string]SECLFunc{
	"nop": func(list *types.MapList) (types.Value, error) {
		return types.Nil{}, nil
	},
}

// EvalMapList executes a MapList which has been marked executable with the correct function
// The first element of the list must be a Function type with a valid identifier
func EvalMapList(list *types.MapList) (types.Value, error) {
	if !list.Executable {
		return nil, ErrNotExecutable
	}
	if len(list.List) == 0 {
		return nil, ErrListLengthZero
	}
	if list.List[0].Type() != types.TFunction {
		return nil, ErrNotAFunction
	}
	fnc := list.List[0].(types.Function)
	fncCall, ok := functions[fnc.Identifier]
	if !ok {
		return nil, ErrUnknownFunction
	}

	return fncCall(list)
}