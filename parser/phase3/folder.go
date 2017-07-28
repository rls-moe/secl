package phase3 // import "go.rls.moe/secl/parser/phase3"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/types"
)

// ErrKeyAtEnd indicates that a key was the last element of a list, which means it has no value to associate with.
var ErrKeyAtEnd = errors.New("Key was at end of list")

type p3Empty struct{}

const tp3Empty types.Type = "P3_EMPTY_REMOVE"

func (*p3Empty) Literal() string {
	return "nil"
}

func (*p3Empty) Type() types.Type {
	return tp3Empty
}

// Fold takes a maplist and folds any mapkeys and their values together
//
// This function will insert placeholder values, p3Empty{}, into any list item it folds into the map
// to make the iteration easier. This means it is necessary to run Clean() afterwards.
func Fold(maplist *types.MapList) error {
	for k := range maplist.List {
		cur := maplist.List[k]
		switch cur.Type() {
		case phase1.TMapExec:
			if len(maplist.List) <= k+1 {
				return ErrKeyAtEnd
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
