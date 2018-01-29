package query // import "go.rls.moe/secl/query"

import (
	"reflect"

	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
)

type PathSegment interface {
	Select(types.Value) (types.Value, error)
}

type Query struct {
	segs []PathSegment
}

func NewQuery(path ...PathSegment) Query {
	return Query{segs: path}
}

func (q Query) Run(list *types.MapList) (types.Value, error) { return q.Select(list) }

func (q Query) Select(val types.Value) (types.Value, error) {
	var curVal = val

	for _, v := range q.segs {
		nextVal, err := v.Select(curVal)
		if err != nil {
			return curVal, err
		}
		if nextVal == nil {
			return nil, errors.Errorf("Query returned empty")
		}
		curVal = nextVal
	}

	return curVal, nil
}

type mapKeySelect struct {
	Key string
}

func KeySelect(key string) PathSegment {
	return &mapKeySelect{Key: key}
}

func (m *mapKeySelect) Select(value types.Value) (types.Value, error) {
	if value.Type() == types.TMapList {
		ml := value.(*types.MapList)
		if ml == nil {
			return nil, errors.New("Map is nil")
		}
		v, ok := ml.Map[types.String{Value: m.Key}]
		if !ok {
			return nil, errors.Errorf("Key not present in map: %s", m.Key)
		}
		return v, nil
	}
	return nil, errors.Errorf("Expected map but got %s for key %s", value.Type(), m.Key)
}

type mapKeySelectDefault struct {
	Key     string
	Default types.Value
}

func KeySelectDefault(key string, def types.Value) PathSegment {
	return &mapKeySelectDefault{Key: key, Default: def}
}

func (m *mapKeySelectDefault) Select(value types.Value) (types.Value, error) {
	if value.Type() == types.TMapList {
		ml := value.(*types.MapList)
		v, ok := ml.Map[types.String{Value: m.Key}]
		if !ok {
			return nil, errors.Errorf("Key not present in map: %s", m.Key)
		}
		return v, nil
	}
	return nil, errors.Errorf("Expected map but got %s for key %s", value.Type(), m.Key)
}

type listSelect struct {
	Index int
}

func ListSelect(index int) PathSegment {
	return &listSelect{Index: index}
}

func (l *listSelect) Select(value types.Value) (types.Value, error) {
	if value.Type() == types.TMapList {
		ml := value.(*types.MapList)
		if l.Index >= len(ml.List) {
			return nil, errors.Errorf("Index exceeds size of map: %d", l.Index)
		}
		return ml.List[l.Index], nil
	}
	return nil, errors.Errorf("Expected map but got %s for key %s", value.Type(), l.Index)
}

type listSelectDefaultStruct struct {
	Index   int
	Default types.Value
}

func ListSelectDefault(index int, def types.Value) PathSegment {
	return &listSelectDefaultStruct{Index: index, Default: def}
}

func (l *listSelectDefaultStruct) Select(value types.Value) (types.Value, error) {
	if value.Type() == types.TMapList {
		ml := value.(*types.MapList)
		if l.Index >= len(ml.List) {
			return l.Default, nil
		}
		return ml.List[l.Index], nil
	}
	return nil, errors.Errorf("Expected map but got %s for key %s", value.Type(), l.Index)
}

type ForceType struct {
	Type types.Type
}

// NewForceType returns a query segment for checking the current value against a certain type
func NewForceType(p types.Type) PathSegment {
	return &ForceType{Type: p}
}

func (f *ForceType) Select(value types.Value) (types.Value, error) {
	if value.Type() == f.Type {
		return value, nil
	}
	return nil, errors.Errorf("Type mismatch of type %s to %s", value.Type(), f.Type)
}

type OnlyList struct{}
type OnlyMap struct{}

func NewOnlyList() PathSegment {
	return OnlyList{}
}

func NewOnlyMap() PathSegment {
	return OnlyMap{}
}

func (OnlyList) Select(value types.Value) (types.Value, error) {
	if value.Type() == types.TMapList {
		return &types.MapList{Map: map[types.String]types.Value{}, List: value.(*types.MapList).List}, nil
	}
	return nil, errors.Errorf("Type mismatch, filtering a list needs a maplist but got a %s", value.Type())
}

func (OnlyMap) Select(value types.Value) (types.Value, error) {
	if value.Type() == types.TMapList {
		return &types.MapList{List: []types.Value{}, Map: value.(*types.MapList).Map}, nil
	}
	return nil, errors.Errorf("Type mismatch, filtering a list needs a maplist but got a %s", value.Type())
}

type Unmarshal struct {
	Target interface{}
}

func NewUnmarshal(target interface{}) Unmarshal {
	return Unmarshal{Target: target}
}

func (u Unmarshal) Select(value types.Value) (types.Value, error) {

	// check if we can unpack directly
	if v, ok := u.Target.(SECLUnmarshal); ok {
		return nil, v.UnmarshalSECL(value)
	}

	if v, ok := u.Target.(PathSegment); ok {
		return v.Select(value)
	}

	// if we can't unpack directly, it must be a pointer value
	rv := reflect.ValueOf(u.Target)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, errors.Errorf("Expected a pointer not a %s for unmarshalling", rv.Kind().String())
	}

	rvp := reflect.Indirect(rv)

	if rvp.Kind() == reflect.String && value.Type() == types.TString {
		rvp.SetString(value.Literal())
		return value, nil
	} else if rvp.Kind() == reflect.Bool && value.Type() == types.TBool {
		rvp.SetBool(value.(*types.Bool).Value)
		return value, nil
	} else if (rvp.Kind() == reflect.Int ||
		rvp.Kind() == reflect.Int64 ||
		rvp.Kind() == reflect.Int32 ||
		rvp.Kind() == reflect.Int16 ||
		rvp.Kind() == reflect.Int8) && value.Type() == types.TInteger {
		rvp.SetInt(value.(*types.Integer).Value.Int64())
		return value, nil
	} else if (rvp.Kind() == reflect.Uint ||
		rvp.Kind() == reflect.Uint64 ||
		rvp.Kind() == reflect.Uint32 ||
		rvp.Kind() == reflect.Uint16 ||
		rvp.Kind() == reflect.Uint8) && value.Type() == types.TInteger {
		rvp.SetUint(value.(*types.Integer).Value.Uint64())
		return value, nil
	} else if (rvp.Kind() == reflect.Float32 ||
		rvp.Kind() == reflect.Float64) && value.Type() == types.TFloat {
		f, _ := value.(*types.Float).Value.Float64()
		rvp.SetFloat(f)
		return value, nil
	}

	return nil, errors.Errorf("Could not unmarshal because of a type mismatch: %s to %s", rvp.Kind(), value.Type())
}

type SECLUnmarshal interface {
	UnmarshalSECL(types.Value) error
}

type UnmarshalQuery struct {
	Target interface{}
	Query  Query
}

func NewUnmarshalWithQuery(target interface{}, path ...PathSegment) UnmarshalQuery {
	return UnmarshalQuery{
		Target: target,
		Query:  NewQuery(path...),
	}
}

func (u UnmarshalQuery) Select(value types.Value) (types.Value, error) {
	v2, err := u.Query.Select(value)
	if err != nil {
		return value, err
	}

	return NewUnmarshal(u.Target).Select(v2)
}

func (u UnmarshalQuery) Run(list *types.MapList) error {
	_, err := u.Select(list)
	return err
}
