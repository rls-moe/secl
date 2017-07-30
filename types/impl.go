package types // import "go.rls.moe/secl/types"

// Literal returns the raw string value of the entity
func (s String) Literal() string {
	return s.Value
}

// Type returns TString
func (String) Type() Type {
	return TString
}

// Literal returns "true" or "false" depending on the internal value
func (b *Bool) Literal() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// Type returns TBool
func (*Bool) Type() Type {
	return TBool
}

// Literal returns the string value of the internal big.Int
func (i *Integer) Literal() string {
	return i.Value.String()
}

// Type returns TInteger
func (*Integer) Type() Type {
	return TInteger
}

// Literal returns the string value of the internal big.Float
func (f *Float) Literal() string {
	return f.Value.String()
}

// Type returns TFloat
func (*Float) Type() Type {
	return TFloat
}

// Literal returns "map:", it does not return the full map contents or it's literals
func (m *MapList) Literal() string {
	return "map:"
}

// Type returns TMapList
func (*MapList) Type() Type {
	return TMapList
}

// Literal returns "func:" and the function identifier
func (f Function) Literal() string {
	return f.Identifier
}

// Type returns TFunction
func (Function) Type() Type {
	return TFunction
}

// Literal returns "nil"
func (Nil) Literal() string {
	return "nil"
}

// Type returns TNil
func (Nil) Type() Type {
	return TNil
}
