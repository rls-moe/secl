package types // import "go.rls.moe/secl/types"

import "math/big"

// Type represents the type of a specific value in the SECL implementation, intermediate stages may also implement their own types for temporary use
type Type string

// These are the basic types of SECL you are able to interact with.
const (
	TMapList  Type = "MAPLIST"
	TString        = "STRING"
	TBool          = "BOOLEAN"
	TInteger       = "NUMBER_INTEGER"
	TFloat         = "NUMBER_FLOAT"
	TFunction      = "FUNCTION"
	TNil           = "NIL"
)

// String is a text value
type String struct {
	Value string
}

var _ Value = &String{} // Assert that String is a Value

// Bool is either true or false
type Bool struct {
	Value bool
}

var _ Value = &Bool{} // Assert that Bool is a Value

// Integer is a big.Int value, meaning arbitrary precision integer. It may exceed 64bit boundaries but may not be usable on all implementations
type Integer struct {
	Value *big.Int
}

var _ Value = &Integer{} // Assert that Integer is a value

// Float is a big.Float value, meaning arbitrary precision float. It may exceed float64 ranges but may not be usable on all implementations
type Float struct {
	Value *big.Float
}

var _ Value = &Float{} // Assert that Float is a Value

// MapList is a combination of maps and lists into one entity. Keys must be a string
type MapList struct {
	Executable bool
	Map        map[String]Value
	List       []Value
}

var _ Value = &MapList{} // Assert that MapList is a Value

// Function is a executable keyword in a maplist that has not been expanded
type Function struct {
	Identifier string
}

var _ Value = &Function{} // Assert that Function is a Value

// Nil represents a null value, it doesn't have any actual value behind it
type Nil struct{}

var _ Value = &Nil{} // Assert that Nil is a Value
