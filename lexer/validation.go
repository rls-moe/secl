package lexer

import "unicode"

func isNotNull(cr rune) bool {
	return cr != 0
}

func isNull(cr rune) bool {
	return cr == 0
}

func isValidChunkRune(cr rune) bool {
	return isNotNull(cr) && isValidInChunk(cr) && !unicode.IsNumber(cr) && !(cr == '!')
}

func isValidInChunk(cr rune) bool {
	return isNotNull(cr) && (!unicode.IsSpace(cr) && !(cr == '@' ||
		cr == ':' ||
		cr == '(' ||
		cr == ')' ||
		cr == '"'))
}

func isValidDigit(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || (cr == '-' ||
		cr == '+'))
}

func isValidInNumber(cr rune) bool {
	return isNotNull(cr) && (unicode.IsNumber(cr) || (cr == '+' ||
		cr == '-' ||
		cr == 'e' ||
		cr == '.' ||
		cr == '*' ||
		cr == '^'))
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
