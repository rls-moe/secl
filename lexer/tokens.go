package lexer // import "go.rls.moe/secl/lexer"
import (
	"github.com/pkg/errors"
	"go.rls.moe/secl/parser/context"
)

// Token is a simple struct to keep track of the position, type and literal of a token
type Token struct {
	Type       context.TokenType
	Literal    string
	Start, End int // Start and End note the Position of a Token, if End=-1 then the Token was not terminated properly, ie a non-terminated string at EOF
}

func RegisterFunctionKeyword(keyword string) error {
	_, ok := context.DefaultKeywords[keyword]
	if ok {
		return errors.New("keyword already used")
	}
	context.DefaultKeywords[keyword] = context.DefaultTokenTypes.Function
	return nil
}
