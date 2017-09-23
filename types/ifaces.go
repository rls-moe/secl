package types // import "go.rls.moe/secl/types"

// Value is a special interface that all types in SECL implement, internal and external ones.
type Value interface {
	// Literal returns a string representation of the internal value
	Literal() string
	// Type returns the data type of the value, use this to determine to which type to cast
	Type() Type
}

type IRandomized interface {
	IsRandom() bool
}

type IPositionInformation interface {
	Position() (int, int)
}

type DebugValue interface {
	DebugPrint() string
}
