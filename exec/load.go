package exec // import "go.rls.moe/secl/exec"
import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/types"
	"io/ioutil"
	"go.rls.moe/secl/parser"
)

// loadv is a SECL Functions that loads a single value from a file, this is done by manually inducing the lexer and phase 1 parser, but not phase 2 and 3
func loadv(list *types.MapList) (types.Value, error) {
	if len(list.List) != 2 {
		return nil, errors.New("loadv requires and permits only 1 parameter")
	}
	filename := list.List[1]

	data, err := ioutil.ReadFile(filename.Literal())
	if err != nil {
		return nil, errors.Wrap(err, "Error reading file for loadv")
	}

	parser := phase1.NewParser(lexer.NewTokenizer(string(data)))
	// Read 1 token
	if err := parser.Step(); err != nil {
		return nil, errors.Wrap(err, "Could not parse value in file")
	}

	if len(parser.FlatAST.FlatNodes) != 1 {
		return nil, errors.Errorf("Parser wanted 1 token but got %d tokens instead", len(parser.FlatAST.FlatNodes))
	}
	tok := parser.FlatAST.FlatNodes[0]

	switch tok.Type() {
	case types.TBool:
		fallthrough
	case types.TString:
		fallthrough
	case types.TFloat:
		fallthrough
	case types.TInteger:
		fallthrough
	case types.TNil:
		fallthrough
	case types.TFunction:
		fallthrough
	case types.TBinary:
		return tok, nil
	default:
		return nil, errors.Errorf("Wanted a token of type bool, string, float, integer, nil, function or binary but got %s instead", tok.Type())
	}
}

func loadb(list *types.MapList) (types.Value, error) {
	if len(list.List) != 2 {
		return nil, errors.New("loadb requires and permits only 1 parameter")
	}

	filename := list.List[1]

	data, err := ioutil.ReadFile(filename.Literal())
	if err != nil {
		return nil, errors.Wrap(err, "Error reading file for loadb")
	}

	return &types.Binary{Raw: data}, nil
}

func subloadfile(filename string) (types.Value, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading file for load*")
	}
	return parser.ParseString(string(data))
}

func loadf(list *types.MapList) (types.Value, error) {
	if len(list.List) != 2 {
		return nil, errors.New("loadf requires and permits only 1 parameter")
	}

	return subloadfile(list.List[1].Literal())
}