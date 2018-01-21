package exec // import "go.rls.moe/secl/exec"

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
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
			val = defVal.Literal()
		}

		if tVal, ok := list.Map[types.String{Value: "type"}]; ok {
			var t types.Type
			switch strings.ToLower(tVal.Literal()) {
			case "integer":
				t = types.TInteger
			case "float":
				fallthrough
			case "number":
				t = types.TFloat
			case "bool":
				t = types.TBool
			case "string":
				t = types.TString
			default:
				return nil, errors.Errorf("type %s for env parameter is unknown", tVal.Literal())
			}
			return types.CoerceType(&types.String{Value: val}, t)
		}

		return &types.String{Value: val}, nil
	})
}
