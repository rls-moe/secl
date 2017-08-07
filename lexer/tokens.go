package lexer // import "go.rls.moe/secl/lexer"
import "github.com/pkg/errors"

// Token is a simple struct to keep track of the position, type and literal of a token
type Token struct {
	Type       TokenType
	Literal    string
	Start, End int // Start and End note the Position of a Token, if End=-1 then the Token was not terminated properly, ie a non-terminated string at EOF
}

// TokenType distinguishes several types of tokens
type TokenType string

// This block declares several Types used by the lexer and the parsers to indicate the nature of a specific token
// They do not indicate validity.
const (
	TTIllegal          TokenType = "ILLEGAL"            // Illegal Character sequences, this should not appear as any character sequence is legal
	TTEOF                        = "EOF"                // EOF indicates either a 0x00 character or an end of file
	TTMapListBegin               = "("                  // Indicates that a new MapList is beginning
	TTMapListEnd                 = ")"                  // Indicates that the current MapList is terminated
	TTModMapKey                  = ":"                  // Mark the previous list entry as a key, use the next as value
	TTModExecMap                 = "!"                  // Execute the following map as function call
	TTModTrim                    = "@"                  // Trim the following string
	TTSingleWordString           = "SINGLE_WORD_STRING" // Any string without spaces
	TTString                     = "STRING"             // Any "" string
	TTNumber                     = "NUMBER"             // 01234567+89.00e-18^19
	TTNil                        = "NIL"                // A nil value, internally equivalent to empty
	TTEmpty                      = "EMPTY"              // In later phases, this is replaced with an empty map ()
	TTBool                       = "BOOLEAN"            // A boolean value
	TTRandstr                    = "RANDSTR"            // A random string value
	TTFunction                   = "FUNCTION"           // A function keyword
	TTComment                    = "COMMENT"            // A comment in the file
)

var keywords = map[string]TokenType{
	"nil":        TTNil,
	"empty":      TTEmpty,
	"nothing":    TTEmpty,
	"true":       TTBool,
	"on":         TTBool,
	"yes":        TTBool,
	"allow":      TTBool,
	"false":      TTBool,
	"off":        TTBool,
	"no":         TTBool,
	"deny":       TTBool,
	"maybe":      TTBool,
	"randstr":    TTRandstr,
	"randstr32":  TTRandstr,
	"randstr64":  TTRandstr,
	"randstr128": TTRandstr,
	"randstr256": TTRandstr,
	"loadd":      TTFunction,
	"loadv":      TTFunction,
	"loadb":      TTFunction,
	"loadf":      TTFunction,
	"decb64":     TTFunction,
	"nop":        TTFunction,
}

func RegisterFunctionKeyword(keyword string) error {
	_, ok := keywords[keyword]
	if ok {
		return errors.New("Keyword already used")
	}
	keywords[keyword] = TTFunction
	return nil
}
