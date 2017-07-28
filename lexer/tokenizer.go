package lexer

import (
	"fmt"
	"unicode"
)

// The below is inspired by https://interpreterbook.com

type Tokenizer struct {
	input     []rune
	pos, rPos int
	curRune   rune
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{
		input: []rune(input),
	}
	t.readChar()
	return t
}

func (t *Tokenizer) readChar() {
	if t.rPos >= len(t.input) {
		t.curRune = 0
	} else {
		t.curRune = t.input[t.rPos]
	}
	t.pos = t.rPos
	t.rPos++
}

func (t *Tokenizer) peekChar() rune {
	if t.rPos >= len(t.input) {
		return 0
	} else {
		return t.input[t.rPos]
	}
}

func (t *Tokenizer) NextToken() Token {
	var tok Token

	t.skipWhitespace()

	switch t.curRune {
	case '(':
		tok = newToken(TT_mapListBegin, t.curRune, t.pos, t.pos)
	case ')':
		tok = newToken(TT_mapListEnd, t.curRune, t.pos, t.pos)
	case '!':
		tok = newToken(TT_mod_execMap, t.curRune, t.pos, t.pos)
	case '@':
		tok = newToken(TT_mod_trim, t.curRune, t.pos, t.pos)
	case ':':
		tok = newToken(TT_mod_mapKey, t.curRune, t.pos, t.pos)
	case 0:
		tok.Literal = ""
		tok.Type = TT_eof
		tok.Start, tok.End = t.pos, t.pos
	default:
		if isValidCommentStarter(t.curRune) {
			tok.Type = TT_comment
			if t.curRune == '/' && t.peekChar() == '*' {
				tok.Literal, tok.Start, tok.End = t.readMlComment()
				t.readChar()
				return tok
			} else if (t.curRune == '/' && t.peekChar() == '/') || t.curRune != '/' {
				tok.Literal, tok.Start, tok.End = t.readSlComment()
				t.readChar()
				return tok
			}
		}
		if isValidDigit(t.curRune) {
			tok.Literal, tok.Start, tok.End = t.readNumber()
			tok.Type = TT_number
			return tok
		} else if isValidChunkRune(t.curRune) {
			tok.Literal, tok.Start, tok.End = t.readChunk()
			if kw, kwtt := resolveKeyword(tok.Literal); kw {
				tok.Type = kwtt
			} else {
				tok.Type = TT_singleWordString
			}
			return tok
		} else if isValidStringBorder(t.curRune) {
			tok.Literal, tok.Start, tok.End = t.readString()
			tok.Type = TT_string
			return tok
		} else {
			fmt.Print("Illegal char at ", t.pos)
			tok = newToken(TT_illegal, t.curRune, t.pos, t.pos)
		}
	}

	t.readChar()
	return tok
}

func (t *Tokenizer) readChunk() (string, int, int) {
	pos := t.pos
	for isValidInChunk(t.curRune) {
		t.readChar()
	}
	return string(t.input[pos:t.pos]), pos, t.pos - 1
}

func (t *Tokenizer) readNumber() (string, int, int) {
	pos := t.pos
	for isValidInNumber(t.curRune) {
		t.readChar()
	}
	return string(t.input[pos:t.pos]), pos, t.pos - 1
}

func (t *Tokenizer) readString() (string, int, int) {
	pos := t.pos
	t.readChar()
	for !isValidStringBorder(t.curRune) {
		t.readChar()
		if isEscapeInString(t.curRune) {
			t.readChar()
			t.readChar()
		}
		if isNull(t.curRune) {
			return string(t.input[pos+1 : t.pos-1]), pos, -1
		}
	}
	t.readChar()
	return string(t.input[pos+1 : t.pos-1]), pos, t.pos - 1
}

func (t *Tokenizer) readMlComment() (string, int, int) {
	pos := t.pos
	for {
		t.readChar()
		if t.curRune == '*' && t.peekChar() == '/' {
			t.readChar()
			t.readChar()
			return string(t.input[pos:t.pos]), pos, t.pos - 1
		}
		if t.curRune == 0 {
			return string(t.input[pos:t.pos]), pos, -1
		}
	}
}

func (t *Tokenizer) readSlComment() (string, int, int) {
	pos := t.pos
	for {
		t.readChar()
		if t.curRune == '\n' {
			return string(t.input[pos:t.pos]), pos, t.pos - 1
		}
	}
}

func (t *Tokenizer) skipWhitespace() {
	for unicode.IsSpace(t.curRune) {
		t.readChar()
	}
}

func newToken(ttype TokenType, cr rune, start, end int) Token {
	return Token{Type: ttype, Literal: string(cr), Start: start, End: end}
}

func resolveKeyword(chunk string) (bool, TokenType) {
	tt, ok := keywords[chunk]
	return ok, tt
}
