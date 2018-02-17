package lexer // import "go.rls.moe/secl/lexer"

import (
	"fmt"
	"unicode"

	"go.rls.moe/secl/parser/context"
)

// The below is inspired by https://interpreterbook.com

// Tokenizer is a simple stream-based/rune-wise lexer. It accepts UTF-8 strings and does not terminate at a 0x00 rune
//
// 0x00 is considered a EOF, however, by calling NextToken you can read beyond any 0x00 character as long as there
// is additional data
//
// Reading beyond the end of the string results in a EOF token
type Tokenizer struct {
	input     []rune
	pos, rPos int
	curRune   rune
	ctx       *context.Lexer
}

// NewTokenizer creates a new Tokenizer based on the given string input
func NewTokenizer(ctx *context.Parser, input string) *Tokenizer {
	t := &Tokenizer{
		input: []rune(input),
		ctx:   ctx.ToLexer(),
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
	}
	return t.input[t.rPos]
}

// NextToken reads the next token from the input stream and returns it.
// If the rune 0x00 is encountered or the tokenizer reads beyond the end of the input, it returns a EOF token
func (t *Tokenizer) NextToken() Token {
	var tok Token

	t.skipWhitespace()

	switch t.curRune {
	case '(':
		tok = newToken(t.ctx.Symbols.MapListBegin, t.curRune, t.pos, t.pos)
	case ')':
		tok = newToken(t.ctx.Symbols.MapListEnd, t.curRune, t.pos, t.pos)
	case '!':
		tok = newToken(t.ctx.Symbols.ModExecMap, t.curRune, t.pos, t.pos)
	case '@':
		tok = newToken(t.ctx.Symbols.ModTrim, t.curRune, t.pos, t.pos)
	case ':':
		tok = newToken(t.ctx.Symbols.ModMapKey, t.curRune, t.pos, t.pos)
	case 0:
		tok.Literal = ""
		tok.Type = t.ctx.Symbols.EOF
		tok.Start, tok.End = t.pos, t.pos
	default:
		if isValidCommentStarter(t.curRune) {
			tok.Type = t.ctx.Symbols.Comment
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
			tok.Type = t.ctx.Symbols.Number
			return tok
		} else if isValidChunkRune(t.curRune) {
			tok.Literal, tok.Start, tok.End = t.readChunk()
			if kw, kwtt := t.resolveKeyword(tok.Literal); kw {
				tok.Type = kwtt
			} else {
				tok.Type = t.ctx.Symbols.SingleWordString
			}
			return tok
		} else if isValidStringBorder(t.curRune) {
			tok.Literal, tok.Start, tok.End = t.readString()
			tok.Type = t.ctx.Symbols.String
			return tok
		}
		fmt.Print("Illegal char at ", t.pos)
		tok = newToken(t.ctx.Symbols.Illegal, t.curRune, t.pos, t.pos)
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
	for isValidInNumber(t.curRune) || (t.pos-pos <= 3 && isValidAsSecondDigit(t.curRune)) {
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

func newToken(ttype context.TokenType, cr rune, start, end int) Token {
	return Token{Type: ttype, Literal: string(cr), Start: start, End: end}
}

func (t *Tokenizer) resolveKeyword(chunk string) (bool, context.TokenType) {
	tt, ok := t.ctx.Keywords[chunk]
	return ok, tt
}
