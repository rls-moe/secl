package phase1

import (
	"go.rls.moe/secl/types"
	"regexp"
	"math/big"
	"github.com/pkg/errors"
)

var (
	regexIntegerDecimal = regexp.MustCompile(`^[+-]?\d+$`)
	regexIntegerHex     = regexp.MustCompile(`^[+-]?0x[0-9A-Fa-f]+$`)
	regexIntegerOct     = regexp.MustCompile(`^[+-]?0o[0-7]+$`)
	regexIntegerBin     = regexp.MustCompile(`^[+-]?0b[01]+$`)
	regexFloatDecimal   = regexp.MustCompile(`^[+-]?\d+.\d+$`)
	regexFloatSci       = regexp.MustCompile(`^[+-]?\d+.\d+\*10\^[+-]?\d{1,3}$`)
	regexFloatExp       = regexp.MustCompile(`^[+-]?\d+.\d+e[+-]?\d{1,3}$`)
)

func ConvertNumber(lit string) (types.Value, error) {
	if regexIntegerDecimal.MatchString(lit) {
		bi := big.NewInt(0)
		bir, ok := bi.SetString(lit, 10)
		if !ok {
			return nil, errors.New("Could not parse integer")
		}
		return &types.Integer{bir}, nil
	} else if regexFloatDecimal.MatchString(lit) {
		bf := big.NewFloat(0)
		bfr, ok := bf.SetString(lit)
		if !ok {
			return nil, errors.New("Could not parse float")
		}
		return &types.Float{bfr}, nil
	}
	return nil, errors.New("Number did not match any known format")
}