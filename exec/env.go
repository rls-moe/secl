package exec // import "go.rls.moe/secl/exec"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
	"os"
)

func init() {
	RegisterFunction("env", func(list *types.MapList) (types.Value, error) {
		if len(list.List) != 2 {
			return nil, errors.New("env needs 1 parameter")
		}

		env := list.List[1]

		if env.Type() != types.TString {
			return nil, errors.New("env expects a string parameter")
		}

		val := os.Getenv(env.Literal())

		defVal, useDefVal := list.Map[types.String{Value: "default"}]
		if val == "" && useDefVal {
			if defVal.Type() != types.TString {
				return nil, errors.New("env default value must be a string")
			}
			val = defVal.Literal()
		}

		return types.String{Value: val}, nil
	})
}
