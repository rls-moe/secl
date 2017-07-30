package exec // import "go.rls.moe/secl/exec"

import (
	"go.rls.moe/secl/types"
	"github.com/pkg/errors"
)

var (
	ErrNoMapList = errors.New("Input was not a maplist")
)
func Eval(value types.Value) (types.Value, error) {
	if value.Type() != types.TMapList {
		return nil, ErrNoMapList
	}

	ml := value.(*types.MapList)

	for k := range ml.Map {
		if ml.Map[k].Type() == types.TMapList {
			if val, err := Eval(ml.Map[k]); err != nil {
				return nil, err
			} else {
				ml.Map[k] = val
			}
		}
	}

	for k := range ml.List {
		if ml.List[k].Type() == types.TMapList {
			if val, err := Eval(ml.List[k]); err != nil {
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