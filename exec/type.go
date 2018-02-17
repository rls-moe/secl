package exec // import "go.rls.moe/secl/exec"

import (
	"strings"

	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/types"
)

const functionNameConvertType = "=>type"
const functionNameConverTypeShort = "=>T"

func init() {
	fun := func(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
		funName := list.List[0].(types.Function).Identifier

		var (
			val        types.Value
			targetType types.Value
		)

		if len(list.List) == 2 {
			val = list.List[1]

			var ok bool
			targetType, ok = list.Map[types.String{Value: "type"}]

			if !ok && len(list.List) == 2 {
				return nil, errors.New(funName + " expects a target type")
			}
		} else if len(list.List) == 3 {
			targetType = list.List[1]
			val = list.List[2]
		} else {
			return nil, errors.New(funName + " needs 1 parameter (type) or two arguments")
		}
		if targetType.Type() != types.TString {
			return nil, errors.New(funName + " target must be a string value")
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
	}
	context.MustRegisterFunction(functionNameConvertType, fun)
	context.MustRegisterFunction(functionNameConverTypeShort, fun)
}
