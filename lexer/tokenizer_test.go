package lexer

import (
	"testing"

	"go.rls.moe/secl/parser/context"
)

func TestTokenizer_NextToken(t *testing.T) {
	input := `( !(nop) let: 909 other:@"test" k2: "test
	hello🤣🤣\"ttt" k3?: k4! decimal: 0.111e-19*10^10 let2: empty let3: nil randstr128 decb64 true off )
# This is a comment
; Also a comment
// Also a comment
/+notacomment
/* this is a ml
comment that
should stop parsing here */ () 0xF 0x0 0b0 0b1 0o7 0o1`

	ctx := context.NewParserContext()
	tests := []struct {
		expectedType    context.TokenType
		expectedLiteral string
		start, end      int
	}{
		{ctx.Symbols.MapListBegin, "(", 0, 0},
		{ctx.Symbols.ModExecMap, "!", 2, 2},
		{ctx.Symbols.MapListBegin, "(", 3, 3},
		{ctx.Symbols.SingleWordString, "nop", 4, 6},
		{ctx.Symbols.MapListEnd, ")", 7, 7},
		{ctx.Symbols.SingleWordString, "let", 9, 11},
		{ctx.Symbols.ModMapKey, ":", 12, 12},
		{ctx.Symbols.Number, "909", 14, 16},
		{ctx.Symbols.SingleWordString, "other", 18, 22},
		{ctx.Symbols.ModMapKey, ":", 23, 23},
		{ctx.Symbols.ModTrim, "@", 24, 24},
		{ctx.Symbols.String, "test", 25, 30},
		{ctx.Symbols.SingleWordString, "k2", 32, 33},
		{ctx.Symbols.ModMapKey, ":", 34, 34},
		{ctx.Symbols.String, "test\n\thello🤣🤣\\\"ttt", 36, 55},
		{ctx.Symbols.SingleWordString, "k3?", 57, 59},
		{ctx.Symbols.ModMapKey, ":", 60, 60},
		{ctx.Symbols.SingleWordString, "k4!", 62, 64},
		{ctx.Symbols.SingleWordString, "decimal", 66, 72},
		{ctx.Symbols.ModMapKey, ":", 73, 73},
		{ctx.Symbols.Number, "0.111e-19*10^10", 75, 89},
		{ctx.Symbols.SingleWordString, "let2", 91, 94},
		{ctx.Symbols.ModMapKey, ":", 95, 95},
		{ctx.Symbols.Empty, "empty", 97, 101},
		{ctx.Symbols.SingleWordString, "let3", 103, 106},
		{ctx.Symbols.ModMapKey, ":", 107, 107},
		{ctx.Symbols.Nil, "nil", 109, 111},
		{ctx.Symbols.Randstr, "randstr128", 113, 122},
		{ctx.Symbols.SingleWordString, "decb64", 124, 129},
		{ctx.Symbols.Bool, "true", 131, 134},
		{ctx.Symbols.Bool, "off", 136, 138},
		{ctx.Symbols.MapListEnd, ")", 140, 140},
		{ctx.Symbols.Comment, "# This is a comment", 142, 160},
		{ctx.Symbols.Comment, "; Also a comment", 162, 177},
		{ctx.Symbols.Comment, "// Also a comment", 179, 195},
		{ctx.Symbols.SingleWordString, "/+notacomment", 197, 209},
		{ctx.Symbols.Comment, "/* this is a ml\ncomment that\nshould stop parsing here */", 211, 266},
		{ctx.Symbols.MapListBegin, "(", 268, 268},
		{ctx.Symbols.MapListEnd, ")", 269, 269},
		{ctx.Symbols.Number, "0xF", 271, 273},
		{ctx.Symbols.Number, "0x0", 275, 277},
		{ctx.Symbols.Number, "0b0", 279, 281},
		{ctx.Symbols.Number, "0b1", 283, 285},
		{ctx.Symbols.Number, "0o7", 287, 289},
		{ctx.Symbols.Number, "0o1", 291, 293},
		{ctx.Symbols.EOF, "", 294, 294},
	}

	l := NewTokenizer(ctx, input)
	//t.Logf("Testing input: %q", input)
	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test %d: literal wrong, expected %+#v (% x) got %+#v (% x)", i, tests[i], tests[i].expectedLiteral, tok, tok.Literal)
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("test %d: type wrong, expected %+v got %+v", i, tests[i], tok)
		}

		if tok.Start != tt.start || tok.End != tt.end {
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i, tests[i].start, tests[i].end, tok.Start, tok.End)
		}

		//t.Logf("Token: %s", tok.Type)
	}
}

func TestTokenizer_NextToken2(t *testing.T) {

	ctx := context.NewParserContext()

	input := `/* abort comment on eof`

	tests := []struct {
		expectedType    context.TokenType
		expectedLiteral string
		start, end      int
	}{
		{ctx.Symbols.Comment, "/* abort comment on eof", 0, -1},
	}

	l := NewTokenizer(ctx, input)
	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test %d: literal wrong, expected %+#v (% x) got %+#v (% x)", i, tests[i], tests[i].expectedLiteral, tok, tok.Literal)
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("test %d: type wrong, expected %+v got %+v", i, tests[i], tok)
		}

		if tok.Start != tt.start || tok.End != tt.end {
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i, tests[i].start, tests[i].end, tok.Start, tok.End)
		}
	}
}

func TestTokenizer_NextToken3(t *testing.T) {

	ctx := context.NewParserContext()

	input := `"abort string on eof`

	tests := []struct {
		expectedType    context.TokenType
		expectedLiteral string
		start, end      int
	}{
		{ctx.Symbols.String, "abort string on eo", 0, -1},
	}

	l := NewTokenizer(ctx, input)
	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test %d: literal wrong, expected %+#v (% x) got %+#v (% x)", i, tests[i], tests[i].expectedLiteral, tok, tok.Literal)
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("test %d: type wrong, expected %+v got %+v", i, tests[i], tok)
		}

		if tok.Start != tt.start || tok.End != tt.end {
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i, tests[i].start, tests[i].end, tok.Start, tok.End)
		}
	}
}

func TestTokenizer_NextToken4(t *testing.T) {

	ctx := context.NewParserContext()

	input := ` `

	tests := []struct {
		expectedType    context.TokenType
		expectedLiteral string
		start, end      int
	}{
		{ctx.Symbols.EOF, "", 1, 1},
	}

	l := NewTokenizer(ctx, input)
	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test %d: literal wrong, expected %+#v (% x) got %+#v (% x)", i, tests[i], tests[i].expectedLiteral, tok, tok.Literal)
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("test %d: type wrong, expected %+v got %+v", i, tests[i], tok)
		}

		if tok.Start != tt.start || tok.End != tt.end {
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i, tests[i].start, tests[i].end, tok.Start, tok.End)
		}
	}
}
