package types

import (
	"math/big"
	"regexp"

	"github.com/pkg/errors"
)

var (
	regexIntegerDecimal = regexp.MustCompile(`^[+-]?\d+$`)
	regexIntegerHex     = regexp.MustCompile(`^[+-]?0x[0-9A-Fa-f]+$`)
	regexIntegerOct     = regexp.MustCompile(`^[+-]?0o[0-7]+$`)
	regexIntegerBin     = regexp.MustCompile(`^[+-]?0b[01]+$`)
	regexFloatDecimal   = regexp.MustCompile(`^[+-]?\d+\.\d+$`)
	regexFloatSci       = regexp.MustCompile(`^([+-]?)(\d+\.\d+)\*10\^([+-]?)(\d{1,3})$`)
	regexFloatExp       = regexp.MustCompile(`^([+-]?)(\d+\.\d+)e([+-]?)(\d{1,3})$`)
)

// ConvertNumber will check the incoming literal against the internal regex list and
// parse it using the correct converted, returning the types.Value interface.
func ConvertNumber(lit string) (Value, error) {
	if regexIntegerDecimal.MatchString(lit) {
		bi := big.NewInt(0)
		bir, ok := bi.SetString(lit, 10)
		if !ok {
			return nil, errors.Errorf("Could not parse integer (dec): %s", lit)
		}
		return &Integer{Value: bir}, nil
	} else if regexFloatDecimal.MatchString(lit) {
		bf := big.NewFloat(0)
		bfr, ok := bf.SetString(lit)
		if !ok {
			return nil, errors.Errorf("Could not parse float (dec): %s", lit)
		}
		return &Float{Value: bfr}, nil
	} else if regexIntegerOct.MatchString(lit) {
		bi := big.NewInt(0)
		bir, ok := bi.SetString(lit[2:], 8)
		if !ok {
			return nil, errors.Errorf("Could not parse integer (oct): %s", lit)
		}
		return &Integer{Value: bir}, nil
	} else if regexIntegerHex.MatchString(lit) {
		bi := big.NewInt(0)
		bir, ok := bi.SetString(lit[2:], 16)
		if !ok {
			return nil, errors.Errorf("Could not parse integer (hex): %s", lit)
		}
		return &Integer{Value: bir}, nil
	} else if regexIntegerBin.MatchString(lit) {
		bi := big.NewInt(0)
		bir, ok := bi.SetString(lit[2:], 2)
		if !ok {
			return nil, errors.Errorf("Could not parse integer (bin): %s", lit)
		}
		return &Integer{Value: bir}, nil
	} else if regexFloatSci.MatchString(lit) {
		matches := regexFloatSci.FindStringSubmatch(lit)
		fpv := big.NewFloat(0.0)
		_, ok := fpv.SetString(matches[2])
		if !ok {
			return nil, errors.Errorf("Could not parse numeric value: %s", lit)
		}
		if matches[1] == "-" {
			fpv.Mul(fpv, big.NewFloat(-1.0))
		}
		fpe := big.NewInt(0.0)
		_, ok = fpe.SetString(matches[4], 10)
		if !ok {
			return nil, errors.Errorf("Could not parse exponent value: %s", lit)
		}
		if matches[3] == "-" {
			fpe.Mul(fpe, big.NewInt(-1))
		}
		fpev := big.NewFloat(1.0)
		for fpe.Uint64() > 0 {
			fpev.Mul(fpev, big.NewFloat(10.0))
			fpe.Sub(fpe, big.NewInt(1))
		}
		return &Float{Value: fpv.Mul(fpv, fpev)}, nil
	} else if regexFloatExp.MatchString(lit) {
		matches := regexFloatExp.FindStringSubmatch(lit)
		fpv := big.NewFloat(0.0)
		_, ok := fpv.SetString(matches[2])
		if !ok {
			return nil, errors.Errorf("Could not parse numeric value: %s", lit)
		}
		if matches[1] == "-" {
			fpv.Mul(fpv, big.NewFloat(-1.0))
		}
		fpe := big.NewInt(0.0)
		_, ok = fpe.SetString(matches[4], 10)
		if !ok {
			return nil, errors.Errorf("Could not parse exponent value: %s", lit)
		}
		if matches[3] == "-" {
			fpe.Mul(fpe, big.NewInt(-1))
		}
		fpev := big.NewFloat(1.0)
		for fpe.Uint64() > 0 {
			fpev.Mul(fpev, big.NewFloat(10.0))
			fpe.Sub(fpe, big.NewInt(1))
		}
		return &Float{Value: fpv.Mul(fpv, fpev)}, nil
	}
	return nil, errors.Errorf("Number did not match any known format: %s", lit)
}
