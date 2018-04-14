package query

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
)

func PathToQuery(path string) ([]PathSegment, error) {
	pathSegs := strings.Split(path, ".")
	queries := []PathSegment{}
	for k := range pathSegs {
		var t PathSegment
		var q = pathSegs[k]
		if q == "#" {
			t = NewOnlyMap()
		} else if q == "[]" {
			t = NewOnlyList()
		} else if strings.HasPrefix(q, "[") {
			i := strings.Trim(q, "[]")
			in, err := strconv.Atoi(i)
			if err != nil {
				return nil, err
			}
			t = ListSelect(in)
		} else {
			t = KeySelect(q)
		}
		queries = append(queries, t)
	}
	return queries, nil
}

// SimpleUnmarshal will use a dot-seperated path to read keys from the map
func SimpleUnmarshal(ml *types.MapList, target interface{}, path string) error {
	queries, err := PathToQuery(path)
	if err != nil {
		return err
	}
	if err := NewUnmarshalWithQuery(target, queries...).Run(ml); err != nil {
		return errors.Wrapf(err, "Could not unwrap '%s'", path)
	}
	return nil
}

// SimpleStructUnmarshal accepts a struct and will unmarshal the maplist using the struct tags
func SimpleStructUnmarshal(v types.Value, target interface{}) error {
	ml, ok := v.(*types.MapList)
	if !ok {
		if unmarshaller, ok := target.(SECLUnmarshal); ok {
			return unmarshaller.UnmarshalSECL(v)
		}
		return errors.New("tried to unpack non-maplist into struct without unmarshaller")
	}
	valueOfTarget := reflect.ValueOf(target).Elem()
	typeOfTarget := reflect.TypeOf(target).Elem()

	for i := 0; i < valueOfTarget.NumField(); i++ {
		field := valueOfTarget.Field(i)
		if !field.CanSet() {
			continue
		}
		typeOfField := typeOfTarget.Field(i)

		path := typeOfField.Tag.Get("secl")
		if path == "" {
			path = typeOfField.Name
		}

		if field.Kind() == reflect.Struct {
			queryPath, err := PathToQuery(path)
			if err != nil {
				return err
			}
			subTarget, err := NewQuery(queryPath...).Run(ml)
			if err != nil {
				return err
			}

			if err := SimpleStructUnmarshal(subTarget, field.Addr().Interface()); err != nil {
				return err
			}
			continue
		} else if field.Kind() == reflect.Array {
			return errors.New("Cannot parse Arrays (yet)")
		} else if field.Kind() == reflect.Slice {
			return errors.New("Cannot parse Slices (yet)")
		}

		if err := SimpleUnmarshal(ml, field.Addr().Interface(), path); err != nil {
			return err
		}
	}
	return nil
}
