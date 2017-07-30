package phase1 // import "go.rls.moe/secl/parser/phase1"

import (
	"go.rls.moe/secl/types"
)

// Several internal and temporary types used by the Phase 1 parser
const (
	TMapBegin types.Type = "P1_MAPBEGIN"
	TMapEnd              = "P1_MAPEND"
	TMapKey              = "P1_MAPKEY"
	TMapEmpty            = "P1_MAPEMPTY"
	TMapExec             = "P1_MAPEXEC"
)

// rootNode is a List-based AST with no depth
type rootNode struct {
	FlatNodes []types.Value
}

// Append puts the given value at the end of the current AST
func (r *rootNode) Append(value types.Value) {
	r.FlatNodes = append(r.FlatNodes, value)
}

// ReplaceLast accepts a replacer function and executes it using the last element as parameter.
// If the given function errors, the replacement is not executed and the error propagates
// If the given function returns a new value, this new value replaces the previously stored one.
func (r *rootNode) ReplaceLast(replacer func(in types.Value) (types.Value, error)) error {
	in := r.FlatNodes[len(r.FlatNodes)-1]
	out, err := replacer(in)
	if err != nil {
		return err
	}
	r.FlatNodes[len(r.FlatNodes)-1] = out
	return nil
}

// MapBegin is a temporary type to indicate the beginning of a map
type MapBegin struct{}

// Literal returns "("
func (MapBegin) Literal() string {
	return string(TMapBegin)
}

// Type returns TMapBegin
func (MapBegin) Type() types.Type {
	return TMapBegin
}

var _ types.Value = MapBegin{} // Assert that MapBegin is a Value

// MapEnd is a temporary type to indicate the end of a map
type MapEnd struct{}

// Literal returns ")"
func (MapEnd) Literal() string {
	return string(TMapEnd)
}

// Type returns TMapEnd
func (MapEnd) Type() types.Type {
	return TMapEnd
}

var _ types.Value = MapEnd{} // Assert that MapEnd is a value

// MapKey is a temporary type used to indicate keys, ie "test:"
type MapKey struct {
	Value string
}

// Literal returns "[<key>]"
func (m MapKey) Literal() string {
	return "[" + m.Value + ":]"
}

// Type returns TMapKey
func (MapKey) Type() types.Type {
	return TMapKey
}

// Key returns the Key of the MapKey as a types.String
func (m MapKey) Key() types.String {
	return types.String{
		Value: m.Value,
	}
}

var _ types.Value = MapKey{} // Assert that MapKey is a value

// EmptyMap is a placeholder where the next phase can insert a proper maplist
type EmptyMap struct{}

// Literal returns "empty"
func (EmptyMap) Literal() string {
	return string(TMapEmpty)
}

// Type returns TMapEmpty
func (EmptyMap) Type() types.Type {
	return TMapEmpty
}

var _ types.Value = EmptyMap{} // Assert that EmptyMap is a Value

// ExecMap indicates the following map is to be marked executable
type ExecMap struct{}

// Literal returns "!exec"
func (ExecMap) Literal() string {
	return string(TMapExec)
}

// Type returns TMapExec
func (ExecMap) Type() types.Type {
	return TMapExec
}

var _ types.Value = ExecMap{} // Assert that ExecMap is a Value
