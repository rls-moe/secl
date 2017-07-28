package phase1 // import "go.rls.moe/secl/parser/phase1"

import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/types"
	"math/big"
	"regexp"
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

// ConvertNumber will check the incoming literal against the internal regex list and
// parse it using the correct converted, returning the types.Value interface.
func ConvertNumber(lit string) (types.Value, error) {
	if regexIntegerDecimal.MatchString(lit) {
		bi := big.NewInt(0)
		bir, ok := bi.SetString(lit, 10)
		if !ok {
			return nil, errors.New("Could not parse integer")
		}
		return &types.Integer{Value: bir}, nil
	} else if regexFloatDecimal.MatchString(lit) {
		bf := big.NewFloat(0)
		bfr, ok := bf.SetString(lit)
		if !ok {
			return nil, errors.New("Could not parse float")
		}
		return &types.Float{Value: bfr}, nil
	} else if regexIntegerHex.MatchString(lit) {
		return nil, errors.New("Hexadecimal Numbers not implemented")
	} else if regexIntegerOct.MatchString(lit) {
		return nil, errors.New("Octal Numbers not implemented")
	} else if regexIntegerBin.MatchString(lit) {
		return nil, errors.New("Binary Numbers not implemented")
	} else if regexFloatSci.MatchString(lit) {
		return nil, errors.New("Scientific Floats not implemented")
	} else if regexFloatExp.MatchString(lit) {
		return nil, errors.New("Exp Floats not implemented")
	}
	return nil, errors.New("Number did not match any known format")
}
