package types

import "math/big"

type Type string

const (
	TMapList Type = "MAPLIST"
	TString       = "STRING"
	TBool         = "BOOLEAN"
	TInteger      = "NUMBER_INTEGER"
	TFloat        = "NUMBER_FLOAT"
	TFunction   = "FUNCTION"
)

type String struct {
	Value string
}

var _ Value = &String{}

type Bool struct {
	Value bool
}

var _ Value = &Bool{}

type Integer struct {
	Value *big.Int
}

var _ Value = &Integer{}

type Float struct {
	Value *big.Float
}

var _ Value = &Float{}

type MapList struct {
	Executable bool
	Map        map[String]Value
	List       []Value
}

var _ Value = &MapList{}

type Function struct {
	Identifier string
}

var _ Value = &Function{}