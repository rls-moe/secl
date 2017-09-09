package exec // import "go.rls.moe/secl/exec"

import (
	"go.rls.moe/secl/types"
	"github.com/pkg/errors"
)

var (
	ErrNoMapList = errors.New("Input was not a maplist")
)
func stepEval(value types.Value) (types.Value, error) {
	if value.Type() != types.TMapList {
		return nil, ErrNoMapList
	}

	ml := value.(*types.MapList)

	for k := range ml.Map {
		if ml.Map[k].Type() == types.TMapList {
			if val, err := stepEval(ml.Map[k]); err != nil {
				return nil, err
			} else {
				ml.Map[k] = val
			}
		}
	}

	for k := range ml.List {
		if ml.List[k].Type() == types.TMapList {
			if val, err := stepEval(ml.List[k]); err != nil {
				return nil, err
			} else {
				ml.List[k] = val
			}
		}
	}

	if ml.Executable {
		return EvalMapList(ml)
	}

	return ml, nil
}

func Eval(value types.Value) (types.Value, error) {
	for canEval(value) {
		newVal, err := stepEval(value)
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