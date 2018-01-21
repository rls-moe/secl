package types

import "github.com/pkg/errors"

// CoerceType will attempt to convert the given value into the given
// type, if it fails, it returns an error. If the given type is equal
// to the values type, the value is returned.
func CoerceType(v Value, t Type) (Value, error) {
	if v.Type() == t {
		return v, nil
	}
	var n Value
	switch t {
	case TNil:
		return Nil{}, nil
	case TFunction:
		return nil, errors.New("Cannot coerce to function")
	case TBinary:
		return nil, errors.New("Cannot coerce to binary data")
	case TBool:
		n = &Bool{}
	case TFloat:
		n = &Float{}
	case TInteger:
		n = &Integer{}
	case TMapList:
		return nil, errors.New("Cannot coerce to maptlist")
	case TString:
		n = &String{}
	default:
		return nil, errors.New("Unknown type in coercion")
	}
	return n, n.FromLiteral(v.Literal())
}
