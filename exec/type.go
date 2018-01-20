package exec // import "go.rls.moe/secl/exec"

import (
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
		var targetVal types.Value
		switch targetTypeS.Value {
		case "string":
			targetVal = new(types.String)
		case "integer":
			targetVal = new(types.Integer)
		case "number":
			targetVal = new(types.Float)
		case "float":
			targetVal = new(types.Float)
		case "boolean":
			targetVal = new(types.Bool)
		default:
			return nil, errors.New(functionNameConvertType + ": unrecognized type value; must be one of 'string, integer, number, boolean'")
		}
		return targetVal, targetVal.FromLiteral(val.Literal())
	})
}
