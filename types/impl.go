package types

func (s *String) Literal() string {
	return s.Value
}

func (*String) Type() Type {
	return TString
}

func (b *Bool) Literal() string {
	if b.Value {
		return "true"
	}
	return "false"
}


func (*Bool) Type() Type {
	return TBool
}

func (i *Integer) Literal() string {
	return i.Value.String()
}

func (*Integer) Type() Type {
	return TInteger
}

func (f *Float) Literal() string {
	return f.Value.String()
}

func (*Float) Type() Type {
	return TFloat
}

func (m *MapList) Literal() string {
	return "map:"
}

func (*MapList) Type() Type {
	return TMapList
}

func (f *Function) Literal() string {
	return "func:" + f.Identifier
}

func (*Function) Type() Type {
	return TFunction
}