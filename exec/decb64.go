package exec // import "go.rls.moe/secl/exec"

import (
	"go.rls.moe/secl/types"
	"github.com/pkg/errors"
	"encoding/base64"
)

func decb64(list *types.MapList) (types.Value, error) {
	var useURL = true
	var useRaw = true

	if v, ok := list.Map[types.String{Value:"urlenc"}]; ok {
		if v.Type() != types.TBool {
			return nil, errors.New("Parameter to urlenc must be boolean")
		}
		useURL = v.(*types.Bool).Value
	}

	if v, ok := list.Map[types.String{Value:"rawenc"}]; ok {
		if v.Type() != types.TBool {
			return nil, errors.New("Parameter to rawenc must be boolean")
		}
	}

	var enc *base64.Encoding
	if useURL && useRaw {
		enc = base64.RawURLEncoding
	} else if useURL && !useRaw {
		enc = base64.URLEncoding
	} else if !useURL && useRaw {
		enc = base64.RawStdEncoding
	} else if !useURL && !useRaw {
		enc = base64.StdEncoding
	}

	if len(list.List) != 2 {
		return nil, errors.New("decb64 only accepts 1 unnamed parameter")
	}

	if list.List[1].Type() != types.TString {
		return nil, errors.New("decb64 only accepts a string parameter for decoding")
	}

	dec, err := enc.DecodeString(list.List[1].Literal())

	if err != nil {
		return nil, errors.Wrap(err, "Could not decode data")
	}

	return &types.Binary{Raw: dec}, nil
}