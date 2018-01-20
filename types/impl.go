package types // import "go.rls.moe/secl/types"
import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"go.rls.moe/secl/helper"
)

// Literal returns the raw string value of the entity
func (s String) Literal() string {
	replacer := strings.NewReplacer(
		"\n", "\\n",
		"\t", "\\t",
		"\"", "\\\"",
		"\\", "\\\\",
	)
	return replacer.Replace(s.Value)
}

// FromLiteral sets the string of the value, it does not return an error
func (s *String) FromLiteral(v string) error {
	replacer := strings.NewReplacer(
		"\\n", "\n",
		"\\t", "\t",
		"\\\"", "\"",
		"\\\\", "\\",
	)
	s.Value = replacer.Replace(v)
	return nil
}

// Type returns TString
func (String) Type() Type {
	return TString
}

// Literal returns "true" or "false" depending on the internal value
func (b *Bool) Literal() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// FromLiteral parses the incoming value as literal boolean
func (b *Bool) FromLiteral(s string) error {
	if s == "true" || s == "on" || s == "allow" || s == "yes" {
		b.Value = true
	} else if s == "false" || s == "off" || s == "deny" || s == "no" {
		b.Value = false
	} else if s == "maybe" {
		b.Value = helper.RndFloat() > 0.501
		b.Randomized = Randomized{Random: true}
	} else {
		return errors.Errorf("Wanted a boolean value but got %q", s)
	}
	return nil
}

// Type returns TBool
func (*Bool) Type() Type {
	return TBool
}

// Literal returns the string value of the internal big.Int
func (i *Integer) Literal() string {
	return i.Value.String()
}

// FromLiteral parses the incoming string as number and fails if the result
// is not an integer
func (i *Integer) FromLiteral(v string) error {
	num, err := ConvertNumber(v)
	if err != nil {
		return err
	}
	switch v := num.(type) {
	case *Integer:
		i.Value = v.Value
	default:
		return errors.New("Parser value was not an integer")
	}
	return nil
}

// Type returns TInteger
func (*Integer) Type() Type {
	return TInteger
}

// Literal returns the string value of the internal big.Float
func (f *Float) Literal() string {
	return f.Value.String()
}

// FromLiteral parses the incoming string as number, it accepts integers
// and floats equally
func (f *Float) FromLiteral(v string) error {
	num, err := ConvertNumber(v)
	if err != nil {
		return err
	}
	switch v := num.(type) {
	case *Integer:
		f.Value = new(big.Float).SetInt(v.Value)
	case *Float:
		f.Value = v.Value
	default:
		return errors.New("Parser value was not an integer or float")
	}
	return nil
}

// Type returns TFloat
func (*Float) Type() Type {
	return TFloat
}

// Literal returns "map:", it does not return the full map contents or it's literals
func (m *MapList) Literal() string {
	return "map:"
}

func (m *MapList) FromLiteral(v string) error {
	return errors.New("Maps cannot be parsed from Literals, yet")
}

// Type returns TMapList
func (*MapList) Type() Type {
	return TMapList
}

// Literal returns "func:" and the function identifier
func (f Function) Literal() string {
	return f.Identifier
}

func (f Function) FromLiteral(v string) error {
	return errors.New("Functions cannot be parsed from Literals")
}

// Type returns TFunction
func (Function) Type() Type {
	return TFunction
}

// Literal returns "nil"
func (Nil) Literal() string {
	return "nil"
}

func (Nil) FromLiteral(v string) error {
	if v == "nil" {
		return nil
	}
	return errors.New("Literal was not nil")
}

// Type returns TNil
func (Nil) Type() Type {
	return TNil
}

func (b *Binary) Literal() string {
	return fmt.Sprintf("!(decb64 \"%s\")", base64.RawURLEncoding.EncodeToString(b.Raw))
}

func (b *Binary) FromLiteral(v string) error {
	if !strings.HasPrefix(v, "!(decb64 \"") || !strings.HasSuffix(v, "\")") {
		return errors.New("Binary must have Functional Prefix")
	}
	v = strings.TrimPrefix(v, "!(decb64 \"")
	v = strings.TrimSuffix(v, "\")")
	bin, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return err
	}
	b.Raw = bin
	return nil
}

func (b *Binary) Type() Type {
	return TBinary
}

func (b *Binary) DebugPrint() string {
	return fmt.Sprintf("0x%X", b.Raw)
}
