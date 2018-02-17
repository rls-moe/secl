package exec // import "go.rls.moe/secl/exec"
import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"go.rls.moe/secl/lexer"
	"go.rls.moe/secl/parser"
	"go.rls.moe/secl/parser/context"
	"go.rls.moe/secl/parser/phase1"
	"go.rls.moe/secl/types"
)

// loadv is a SECL Functions that loads a single value from a file, this is done by manually inducing the lexer and phase 1 parser, but not phase 2 and 3
func loadv(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
	if len(list.List) != 2 {
		return nil, errors.New("loadv requires and permits only 1 parameter")
	}
	filename := list.List[1]

	data, err := ioutil.ReadFile(filename.Literal())
	if err != nil {
		return nil, errors.Wrap(err, "Error reading file for loadv")
	}
	defaultCtx := context.NewParserContext()
	parser := phase1.NewParser(defaultCtx, lexer.NewTokenizer(defaultCtx, string(data)))
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

func loadb(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
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
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// If the file does not exist, return an empty map
		return &types.MapList{
			Map:  map[types.String]types.Value{},
			List: []types.Value{},
		}, nil
	}
	if stat.IsDir() {
		return nil, errors.New("file was a directory")
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading file for load*")
	}
	return parser.ParseString(string(data))
}

func loadf(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
	if len(list.List) != 2 {
		return nil, errors.New("loadf requires and permits only 1 parameter")
	}

	return subloadfile(list.List[1].Literal())
}

func loadd(ctx *context.Runtime, list *types.MapList) (types.Value, error) {
	folder, ok := list.Map[types.String{Value: "dir"}]
	if !ok {
		folder = &types.String{Value: "./conf.d"}
	}
	if folder.Type() != types.TString {
		return nil, errors.New("Folder must be a string")
	}
	suffix, ok := list.Map[types.String{Value: "suffix"}]
	if !ok {
		folder = &types.String{Value: ".secl"}
	}
	if suffix.Type() != types.TString {
		return nil, errors.New("Suffix must be a string")
	}
	var loader string
	loaderML, ok := list.Map[types.String{Value: "parser"}]
	if !ok {
		loader = "secl"
	} else {
		loaderI, ok := loaderML.(*types.String)
		if !ok {
			return nil, errors.New("parser parameter must be a string")
		}
		loader = loaderI.Value
	}
	loader = strings.ToLower(loader)
	files, err := filepath.Glob(filepath.Join(folder.(*types.String).Value, "*"+suffix.(*types.String).Value))
	if err != nil {
		return nil, err
	}
	if loader == "plain" {
		var out = &types.MapList{
			Map:  map[types.String]types.Value{},
			List: []types.Value{},
		}
		for k := range files {
			stat, err := os.Stat(files[k])
			if err != nil {
				return nil, err
			}
			if stat.IsDir() {
				return nil, errors.New("plain parser does not read directories")
			}
			cnt, err := ioutil.ReadFile(files[k])
			if err != nil {
				return nil, err
			}
			out.Map[types.String{Value: files[k]}] = &types.String{
				Value: string(cnt),
			}
		}
		return out, nil
	}
	var out = &types.MapList{
		Map: map[types.String]types.Value{},
		List: []types.Value{
			types.Function{Identifier: "merge"},
		},
		Executable: true,
	}
	for k := range files {
		if loader == "secl" {
			ml, err := subloadfile(files[k])
			if err != nil {
				return nil, errors.Wrapf(err, "Could not load file %s", files[k])
			}
			out.List = append(out.List, ml)
		} else if f, ok := ctx.Functions["__loadd_parser@"+loader]; ok {
			l, err := f(ctx, &types.MapList{
				List: []types.Value{
					&types.String{Value: files[k]},
				},
			})
			if err != nil {
				return nil, err
			}
			out.List = append(out.List, l)
		} else {
			return nil, errors.New("Unknown parser " + loader)
		}
	}
	return out, nil
}
