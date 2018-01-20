package types // import "go.rls.moe/secl/types"
import (
	"fmt"
	"sort"
	"strings"
)

// PrintValue will print out a debug string of a given Value. This string is not necessarily a valid SECL
func PrintValue(p Value) string {
	switch p.Type() {
	case TMapList:
		q := p.(*MapList)
		var out = "( "

		if q.Executable {
			out = "!" + out
		}

		var keys = make([]string, 0)
		for k := range q.Map {
			keys = append(keys, k.Value)
		}

		sort.Strings(keys)

		for k := range keys {
			out += PrintValue(&String{Value: keys[k]}) + ": " + PrintValue(q.Map[String{Value: keys[k]}]) + " "
		}

		for k := range q.List {
			out += PrintValue(q.List[k]) + " "
		}

		return out + ")"
	case TString:
		return fmt.Sprintf("\"%s\"", p.Literal())
	default:
		return fmt.Sprintf("%s", p.Literal())
	}
}

// PrintReproducableValue is intended for serializing maplists into a string that reproduces the exact same maplist when read
func PrintReproducableValue(p *MapList) string {
	str := PrintValue(p)
	str = strings.TrimSuffix(str, " )")
	str = strings.TrimPrefix(str, "( ")
	return str
}

// PrintDebug will print out a debug string of a given Value. This contains more information than PrintValue
func PrintDebug(p Value) string {
	switch p.Type() {
	case TMapList:
		q := p.(*MapList)
		var out = "( "

		if q.Executable {
			out = "exec:" + out
		}

		var keys = make([]string, 0)
		for k := range q.Map {
			keys = append(keys, k.Value)
		}

		sort.Strings(keys)

		out += "//MAP "

		for k := range keys {
			out += PrintDebug(&String{Value: keys[k]}) + ": " + PrintDebug(q.Map[String{Value: keys[k]}]) + " "
		}

		out += "//LIST "

		for k := range q.List {
			out += PrintDebug(q.List[k]) + " "
		}

		return out + ")"
	default:
		var str string
		if v, ok := p.(IRandomized); ok && v.IsRandom() {
			str += "{RANDOM}"
		} else {
			if _, ok := p.(*String); ok {
				str += fmt.Sprintf("\"%s\"", p.Literal())
			} else if v, ok := p.(DebugValue); ok {
				str += v.DebugPrint()
			} else {
				str += p.Literal()
			}
		}
		if p.Literal() == string(p.Type()) {
			str += "//"
		} else {
			str += "/" + string(p.Type())
		}
		if v, ok := p.(IPositionInformation); ok {
			s, e := v.Position()
			str += fmt.Sprintf("(%d:%d)", s, e)
		}
		return str
	}
}
