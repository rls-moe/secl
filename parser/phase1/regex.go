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
	regexFloatSci       = regexp.MustCompile(`^([+-]?)(\d+.\d+)\*10\^([+-]?)(\d{1,3})$`)
	regexFloatExp       = regexp.MustCompile(`^([+-]?)(\d+.\d+)e([+-]?)(\d{1,3})$`)
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
		matches := regexFloatSci.FindStringSubmatch(lit)
		fpv := big.NewFloat(0.0)
		_, ok := fpv.SetString(matches[2])
		if !ok {
			return nil, errors.New("Could not parse Numeric Value")
		}
		if matches[1] == "-" {
			fpv.Mul(fpv, big.NewFloat(-1.0))
		}
		fpe := big.NewInt(0.0)
		_, ok = fpe.SetString(matches[4], 10)
		if !ok {
			return nil, errors.New("Could not parse Exponent Value")
		}
		if matches[3] == "-" {
			fpe.Mul(fpe, big.NewInt(-1))
		}
		fpev := big.NewFloat(1.0)
		for fpe.Uint64() > 0 {
			fpev.Mul(fpev, big.NewFloat(10.0))
			fpe.Sub(fpe, big.NewInt(1))
		}
		return &types.Float{Value:fpv.Mul(fpv, fpev)}, nil
	} else if regexFloatExp.MatchString(lit) {
		matches := regexFloatExp.FindStringSubmatch(lit)
		fpv := big.NewFloat(0.0)
		_, ok := fpv.SetString(matches[2])
		if !ok {
			return nil, errors.New("Could not parse Numeric Value")
		}
		if matches[1] == "-" {
			fpv.Mul(fpv, big.NewFloat(-1.0))
		}
		fpe := big.NewInt(0.0)
		_, ok = fpe.SetString(matches[4], 10)
		if !ok {
			return nil, errors.New("Could not parse Exponent Value")
		}
		if matches[3] == "-" {
			fpe.Mul(fpe, big.NewInt(-1))
		}
		fpev := big.NewFloat(1.0)
		for fpe.Uint64() > 0 {
			fpev.Mul(fpev, big.NewFloat(10.0))
			fpe.Sub(fpe, big.NewInt(1))
		}
		return &types.Float{Value:fpv.Mul(fpv, fpev)}, nil
	}
	return nil, errors.New("Number did not match any known format")
}
