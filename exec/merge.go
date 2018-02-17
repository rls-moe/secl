package exec // import "go.rls.moe/secl/exec"
import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/types"
)

func merge(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
	if len(list.Map) != 0 {
		return nil, errors.New("merge does not accept map-key values as parameter")
	}
	if len(list.List) == 1 {
		// Running merge on an empty list results in a nil value
		return types.Nil{}, nil
	}

	inputmaps := list.List[1:]

	var outmap = &types.MapList{
		Map:  map[types.String]types.Value{},
		List: []types.Value{},
	}

	for k := range inputmaps {
		val := inputmaps[k]
		if val.Type() != types.TMapList {
			return nil, errors.New("merge does not accept non-maps in it's primary list")
		}
		if vmap, ok := val.(*types.MapList); ok {
			if err := mergeIntoMap(outmap, vmap); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("merge encountered a non-map in the primary list which lied about it's type")
		}
	}
	return outmap, nil
}

func mergeIntoMap(prim, in *types.MapList) error {
	// Merge maps
	for k := range in.Map {
		if _, ok := prim.Map[k]; !ok {
			prim.Map[k] = in.Map[k]
		} else {
			if in.Map[k].Type() == types.TMapList {
				if prim.Map[k].Type() != types.TMapList {
					return errors.New("Attempted to merge maplist into non-maplist")
				} else {
					// If both are a map, recurse into them
					if err := mergeIntoMap(prim.Map[k].(*types.MapList), in.Map[k].(*types.MapList)); err != nil {
						return err
					}
				}
			} else {
				if prim.Map[k].Type() != in.Map[k].Type() {
					return errors.New("Attempted to change type of existing value")
				} else {
					prim.Map[k] = in.Map[k]
				}
			}
		}
	}

	// Append to lists
	for k := range in.List {
		prim.List = append(prim.List, in.List[k])
	}
	return nil
}
