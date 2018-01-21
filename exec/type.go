package exec // import "go.rls.moe/secl/exec"

import (
	"strings"

	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
)

const functionNameConvertType = "=>type"

func init() {
	RegisterFunction(functionNameConvertType, func(list *types.MapList) (types.Value, error) {
		if len(list.List) != 2 {
			return nil, errors.New(functionNameConvertType + " needs 1 parameter")
		}

		val := list.List[1]

		targetType, ok := list.Map[types.String{Value: "type"}]

		if !ok {
			return nil, errors.New(functionNameConvertType + " expects a target type")
		}
		if targetType.Type() != types.TString {
			return nil, errors.New(functionNameConvertType + " target must be a string value")
		}
		targetTypeS := targetType.(*types.String)
		var targetTypeT types.Type
		switch strings.ToLower(targetTypeS.Value) {
		case "string":
			targetTypeT = types.TString
		case "integer":
			targetTypeT = types.TInteger
		case "float":
			fallthrough
		case "number":
			targetTypeT = types.TFloat
		case "bool":
			targetTypeT = types.TBool
		}
		return types.CoerceType(val, targetTypeT)
	})
}
