package exec // import "go.rls.moe/secl/exec"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/types"
)

var (
	ErrNoMapList = errors.New("Input was not a maplist")
)

func stepEval(ctx *context.Runtime, value types.Value) (types.Value, error) {
	if value.Type() != types.TMapList {
		return nil, ErrNoMapList
	}

	ml := value.(*types.MapList)

	if !ml.InterruptExec {
		for k := range ml.Map {
			if ml.Map[k].Type() == types.TMapList {
				if val, err := stepEval(ctx, ml.Map[k]); err != nil {
					return nil, err
				} else {
					ml.Map[k] = val
				}
			}
		}

		for k := range ml.List {
			if ml.List[k].Type() == types.TMapList {
				if val, err := stepEval(ctx, ml.List[k]); err != nil {
					return nil, err
				} else {
					ml.List[k] = val
				}
			}
		}
	}

	{
		artFML := &types.MapList{
			List: []types.Value{ml},
			Map:  map[types.String]types.Value{},
		}
		for k := range ml.List {
			switch v := ml.List[k].(type) {
			case *types.MapList:
				if v.MergeUp {
					artFML.List = append(artFML.List, v)
				}
			default:
				continue
			}
		}
		if len(artFML.List) > 1 {
			newML, err := merge(ctx, artFML)
			if err != nil {
				return nil, err
			}
			ml = newML.(*types.MapList)
		}
	}

	if ml.Executable {
		return EvalMapList(ctx, ml)
	}

	return ml, nil
}

func Eval(ctx *context.Runtime, value types.Value) (types.Value, error) {
	for canEval(value) {
		newVal, err := stepEval(ctx, value)
		if err != nil {
			return value, err
		}
		value = newVal
	}
	return value, nil
}

func canEval(value types.Value) bool {
	if value.Type() != types.TMapList {
		return false
	}

	ml := value.(*types.MapList)

	if ml.Executable {
		return true
	}

	if ml.InterruptExec {
		return false
	}

	for k := range ml.Map {
		if ml.Map[k].Type() == types.TMapList {
			if ml.Map[k].(*types.MapList).Executable {
				return true
			}
			if canEval(ml.Map[k]) {
				return true
			}
		}
	}

	for k := range ml.List {
		if ml.List[k].Type() == types.TMapList {
			if ml.List[k].(*types.MapList).Executable {
				return true
			}
			if canEval(ml.List[k]) {
				return true
			}
		}
	}

	return false
}
