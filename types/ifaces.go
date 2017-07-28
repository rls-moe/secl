package types

type Value interface {
	Literal() string
	Type() Type
}