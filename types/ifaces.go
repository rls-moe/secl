package types // import "go.rls.moe/secl/types"

// Value is a special interface that all types in SECL implement, internal and external ones.
type Value interface {
	// Literal returns a string representation of the internal value
	Literal() string
	// FromLiteral parses the given string into the Value or returns an error
	// when an error is returned, the original value of the Value must be preserved
	FromLiteral(string) error
	// Type returns the data type of the value, use this to determine to which type to cast
	Type() Type
}

// IRandomized indicates the underlying type contains random data and any debug
// output should indicate this and not produce the actual random value
type IRandomized interface {
	IsRandom() bool
}

// DebugValue indicates the type has it's own debug output function and the default
// string literal function should not be used
type DebugValue interface {
	DebugPrint() string
}
