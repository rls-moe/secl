package context

import "go.rls.moe/secl/types"

// SECLFunc is a generic function to be executed in a SECL file
// It receives the maplist that contained the function
// When it returns no error, it must return a non-nil types.Value entity that
// replaces/expands the function position
type SECLFunc func(ctx *Runtime, list *types.MapList) (types.Value, error)

type TokenType string

// TokenType distinguishes several types of tokens
type TokenTypeRegistry struct {
	Illegal          TokenType
	EOF              TokenType
	MapListBegin     TokenType
	MapListEnd       TokenType
	ModMapKey        TokenType
	ModExecMap       TokenType
	ModTrim          TokenType
	SingleWordString TokenType
	String           TokenType
	Number           TokenType
	Nil              TokenType
	Empty            TokenType
	Bool             TokenType
	Randstr          TokenType
	Function         TokenType
	Comment          TokenType
}

func GetDefaultTokenTypes() *TokenTypeRegistry {
	return &TokenTypeRegistry{
		Illegal:          "ILLEGAL",
		EOF:              "EOF",
		MapListBegin:     "(",
		MapListEnd:       ")",
		ModMapKey:        ":",
		ModExecMap:       "!",
		ModTrim:          "@",
		SingleWordString: "SINGLE_WORD_STRING",
		String:           "STRING",
		Number:           "NUMBER",
		Nil:              "NIL",
		Empty:            "EMPTY",
		Bool:             "BOOLEAN",
		Randstr:          "RANDSTR",
		Function:         "FUNCTION",
		Comment:          "COMMENT",
	}
}
