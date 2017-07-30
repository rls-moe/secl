package phase3 // import "go.rls.moe/secl/parser/phase3"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/types"
)

// ErrKeyAtEnd indicates that a key was the last element of a list, which means it has no value to associate with.
var ErrKeyAtEnd = errors.New("Key was at end of list")

// ErrExecNotMap indicates that a executable marker was followed by something that isn't a map
var ErrExecNotMap = errors.New("Tried to execute non-map")

type p3Empty struct{}

const tp3Empty types.Type = "P3_EMPTY_REMOVE"

func (*p3Empty) Literal() string {
	return string(tp3Empty)
}

func (*p3Empty) Type() types.Type {
	return tp3Empty
}

// Fold takes a maplist and folds any mapkeys and their values together
//
// This function will insert placeholder values, p3Empty{}, into any list item it folds into the map
// to make the iteration easier. This means it is necessary to run Clean() afterwards.
func Fold(maplist *types.MapList) error {
	for k := 0; k < len(maplist.List); k++ {
		cur := maplist.List[k]
		switch cur.Type() {
		case phase1.TMapExec:
			if len(maplist.List) <= k+1 {
				return ErrKeyAtEnd
			}
			if maplist.List[k+1].Type() != types.TMapList {
				return ErrExecNotMap
			}
			maplist.List[k+1].(*types.MapList).Executable = true
			maplist.List[k] = &p3Empty{}
		case types.TMapList:
			mpl := cur.(*types.MapList)
			if err := Fold(mpl); err != nil {
				return err
			}
			maplist.List[k] = mpl
		case phase1.TMapKey:
			if len(maplist.List) <= k+1 {
				return ErrKeyAtEnd
			}
			nxt := maplist.List[k+1]
			if nxt.Type() == types.TMapList {
				mplnxt := nxt.(*types.MapList)
				if err := Fold(mplnxt); err != nil {
					return err
				}
				nxt = mplnxt
			}
			if nxt.Type() == phase1.TMapExec {
				if len(maplist.List) <= k+2 {
					return ErrKeyAtEnd
				}
				nxtnxt := maplist.List[k+2]
				if nxtnxt.Type() != types.TMapList {
					return ErrExecNotMap
				}
				if err := Fold(nxtnxt.(*types.MapList)); err != nil {
					return err
				}
				nxtnxt.(*types.MapList).Executable = true
				nxt = nxtnxt
				maplist.List[k+2] = &p3Empty{}
			}
			mpk := cur.(*phase1.MapKey).Key()
			maplist.Map[mpk] = nxt
			maplist.List[k+1] = &p3Empty{}
			maplist.List[k] = &p3Empty{}
		default:
			continue
		}
	}
	return nil
}

// Clean removes any placeholders put into place by Fold
func Clean(list *types.MapList) {
	for k := 0; k < len(list.List); k++ {
		if list.List[k].Type() == tp3Empty {
			if k < len(list.List)-1 {
				list.List = append(list.List[:k], list.List[k+1:]...)
			} else {
				list.List = list.List[:k]
			}
			k--
		} else if list.List[k].Type() == types.TMapList {
			Clean(list.List[k].(*types.MapList))
		}
	}

	keys := make([]types.String, len(list.Map))
	i := 0
	for k := range list.Map {
		keys[i] = k
		i++
	}

	for j := 0; j < len(keys); j++ {
		if list.Map[keys[j]].Type() == tp3Empty {
			delete(list.Map, keys[j])
		}

		if list.Map[keys[j]].Type() == types.TMapList {
			Clean(list.Map[keys[j]].(*types.MapList))
		}
	}
}
