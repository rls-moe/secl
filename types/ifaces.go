package types // import "go.rls.moe/secl/types"

// Value is a special interface that all types in SECL implement, internal and external ones.
type Value interface {
	Literal() string
	Type() Type
}
