package lexer // import "go.rls.moe/secl/lexer"

import "unicode"

// isNotNull checks if the given rune is equal to 0x00
func isNotNull(cr rune) bool {
	return cr != 0
}

// isNull is the exact inverse of isNotNull
func isNull(cr rune) bool {
	return cr == 0
}

// isValidChunkRune checks if a given character can begin a chunk
func isValidChunkRune(cr rune) bool {
	return isNotNull(cr) && isValidInChunk(cr) && !unicode.IsNumber(cr) && !(cr == '!')
}

// isValidInChunk checks if a given character is valid to be used within a chunk
func isValidInChunk(cr rune) bool {
	return isNotNull(cr) && (!unicode.IsSpace(cr) && !(cr == '@' ||
		cr == ':' ||
		cr == '(' ||
		cr == ')' ||
		cr == '"'))
}

// isValidDigit checks if a given character can start a number token
func isValidDigit(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || (cr == '-' ||
		cr == '+'))
}

// isValidInNumber checks if a given character is valid within a number
func isValidInNumber(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || (cr == '+' ||
		cr == '-' ||
		cr == 'e' ||
		cr == '.' ||
		cr == '*' ||
		cr == '^' ||
		unicode.In(cr, unicode.Hex_Digit)))
}

func isValidAsSecondDigit(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || isValidInNumber(cr) || (cr == 'x' ||
		cr == 'o' ||
		cr == 'b' ||
		cr == 'X' ||
		cr == 'O' ||
		cr == 'B'))
}

func isValidStringBorder(cr rune) bool {
	return cr == '"'
}

func isEscapeInString(cr rune) bool {
	return cr == '\\'
}

func isValidCommentStarter(cr rune) bool {
	return cr == '/' || cr == ';' || cr == '#'
}
