package query

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
	"reflect"
	"strings"
)

func PathToQuery(path string) []PathSegment {
	pathSegs := strings.Split(path, ".")
	queries := []PathSegment{}
	for k := range pathSegs {
		queries = append(queries, KeySelect(pathSegs[k]))
	}
	return queries
}

// SimpleUnmarshal will use a dot-seperated path to read keys from the map
func SimpleUnmarshal(ml *types.MapList, target interface{}, path string) error {
	queries := PathToQuery(path)
	if err := NewUnmarshalWithQuery(target, queries...).Run(ml); err != nil {
		return errors.Wrapf(err, "Could not unwrap %s", path)
	}
	return nil
}

// SimpleStructUnmarshal accepts a struct and will unmarshal the maplist using the struct tags
func SimpleStructUnmarshal(ml *types.MapList, target interface{}) error {
	if unmarshaller, ok := target.(SECLUnmarshal); ok {
		return unmarshaller.UnmarshalSECL(ml)
	}
	valueOfTarget := reflect.ValueOf(target).Elem()
	typeOfTarget := reflect.TypeOf(target).Elem()

	for i := 0; i < valueOfTarget.NumField(); i++ {
		field := valueOfTarget.Field(i)
		typeOfField := typeOfTarget.Field(i)

		path := typeOfField.Tag.Get("secl")
		if path == "" {
			path = typeOfField.Name
		}
		if field.Kind() == reflect.Struct {
			queryPath := PathToQuery(path)
			subTarget, err := NewQuery(queryPath...).Run(ml)
			if err != nil {
				return err
			}
			mlSubTarget, ok := subTarget.(*types.MapList)
			if !ok {
				return errors.New("Struct Subtarget is not a map")
			}
			if err := SimpleStructUnmarshal(mlSubTarget, field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}
		if field.Kind() == reflect.Array || field.Kind() == reflect.Slice || field.Kind() == reflect.Map {
			return errors.New("Cannot parse Maps, Arrays or Slices (yet)")
		}
		if err := SimpleUnmarshal(ml, field.Addr().Interface(), path); err != nil {
			return err
		}
	}
	return nil
}
