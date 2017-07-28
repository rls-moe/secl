package types // import "go.rls.moe/secl/types"

// PrintValue will print out a debug string of a given Value. This string is not parsable but usuable for
// human consumption.
func PrintValue(p Value) string {
	switch p.Type() {
	case TMapList:
		q := p.(*MapList)
		var out = "( "

		if q.Executable {
			out = "exec:" + out
		}

		for k := range q.Map {
			out += PrintValue(&k) + ": " + PrintValue(q.Map[k]) + " "
		}

		for k := range q.List {
			out += PrintValue(q.List[k]) + " "
		}

		return out + ")"
	case TString:
		return "\"" + p.Literal() + "\""
	default:
		return p.Literal()
	}
}
