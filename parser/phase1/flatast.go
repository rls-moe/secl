package phase1

import (
	"go.rls.moe/secl/types"
)

const (
	TMapBegin types.Type = "P1_MAPBEGIN"
	TMapEnd = "P1_MAPEND"
	TMapKey = "P1_MAPKEY"
	TMapEmpty = "P1_MAPEMPTY"
	TMapExec = "P1_MAPEXEC"
)

type RootNode struct {
	FlatNodes []types.Value
}

func (r *RootNode) Append(value types.Value) {
	r.FlatNodes = append(r.FlatNodes, value)
}

func (r *RootNode) ReplaceLast(replacer func(in types.Value) (types.Value, error)) error {
	in := r.FlatNodes[len(r.FlatNodes)-1]
	out, err := replacer(in)
	if err != nil {
		return err
	}
	r.FlatNodes[len(r.FlatNodes)-1] = out
	return nil
}

type MapBegin struct {}

func (*MapBegin) Literal() string {
	return "("
}

func (*MapBegin) Type() types.Type {
	return TMapBegin
}

var _ types.Value = &MapBegin{}

type MapEnd struct {}

func (*MapEnd) Literal() string {
	return ")"
}

func (*MapEnd) Type() types.Type {
	return TMapEnd
}

var _ types.Value = &MapEnd{}

type MapKey struct {
	Value string
}

func (m *MapKey) Literal() string {
	return "[" + m.Value + ":]"
}

func (*MapKey) Type() types.Type {
	return TMapKey
}

func (m *MapKey) Key() types.String {
	return types.String{
		Value: m.Value,
	}
}

type EmptyMap struct {}

func (*EmptyMap) Literal() string {
	return "empty"
}

func (*EmptyMap) Type() types.Type {
	return TMapEmpty
}

type ExecMap struct{}

func (*ExecMap) Literal() string {
	return "!exec"
}

func (*ExecMap) Type() types.Type {
	return TMapExec
}


