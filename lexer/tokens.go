package lexer

type Token struct {
	Type       TokenType
	Literal    string
	Start, End int // Start and End note the Position of a Token, if End=-1 then the Token was not terminated properly, ie a non-terminated string at EOF
}

type TokenType string

const (
	TT_illegal          TokenType = "ILLEGAL"
	TT_eof                        = "EOF"
	TT_mapListBegin               = "("
	TT_mapListEnd                 = ")"
	TT_mod_mapKey                 = ":"
	TT_mod_execMap                = "!"
	TT_mod_trim                   = "@"
	TT_singleWordString           = "SINGLE_WORD_STRING" // Any string without spaces
	TT_string                     = "STRING"             // Any "" string
	TT_number                     = "NUMBER"             // 01234567+89.00e-18^19
	TT_nil                        = "NIL"
	TT_empty                      = "EMPTY"
	TT_bool                       = "BOOLEAN"
	TT_randstr                    = "RANDSTR"
	TT_function                   = "FUNCTION"
	TT_comment                    = "COMMENT"
)

var keywords = map[string]TokenType{
	"nil":        TT_nil,
	"empty":      TT_empty,
	"nothing":    TT_empty,
	"true":       TT_bool,
	"on":         TT_bool,
	"yes":        TT_bool,
	"allow":      TT_bool,
	"false":      TT_bool,
	"off":        TT_bool,
	"no":         TT_bool,
	"deny":       TT_bool,
	"randstr":    TT_randstr,
	"randstr32":  TT_randstr,
	"randstr64":  TT_randstr,
	"randstr128": TT_randstr,
	"randstr256": TT_randstr,
	"loadd":      TT_function,
	"loadv":      TT_function,
	"loadb":      TT_function,
	"decb64":     TT_function,
	"nop":        TT_function,
}
