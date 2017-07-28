package lexer

import (
	"testing"
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
should stop parsing here */ ()`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
		start, end      int
	}{
		{TTMapListBegin, "(", 0, 0},
		{TTModExecMap, "!", 2, 2},
		{TTMapListBegin, "(", 3, 3},
		{TTFunction, "nop", 4, 6},
		{TTMapListEnd, ")", 7, 7},
		{TTSingleWordString, "let", 9, 11},
		{TTModMapKey, ":", 12, 12},
		{TTNumber, "909", 14, 16},
		{TTSingleWordString, "other", 18, 22},
		{TTModMapKey, ":", 23, 23},
		{TTModTrim, "@", 24, 24},
		{TTString, "test", 25, 30},
		{TTSingleWordString, "k2", 32, 33},
		{TTModMapKey, ":", 34, 34},
		{TTString, "test\n\thello🤣🤣\\\"ttt", 36, 55},
		{TTSingleWordString, "k3?", 57, 59},
		{TTModMapKey, ":", 60, 60},
		{TTSingleWordString, "k4!", 62, 64},
		{TTSingleWordString, "decimal", 66, 72},
		{TTModMapKey, ":", 73, 73},
		{TTNumber, "0.111e-19*10^10", 75, 89},
		{TTSingleWordString, "let2", 91, 94},
		{TTModMapKey, ":", 95, 95},
		{TTEmpty, "empty", 97, 101},
		{TTSingleWordString, "let3", 103, 106},
		{TTModMapKey, ":", 107, 107},
		{TTNil, "nil", 109, 111},
		{TTRandstr, "randstr128", 113, 122},
		{TTFunction, "decb64", 124, 129},
		{TTBool, "true", 131, 134},
		{TTBool, "off", 136, 138},
		{TTMapListEnd, ")", 140, 140},
		{TTComment, "# This is a comment", 142, 160},
		{TTComment, "; Also a comment", 162, 177},
		{TTComment, "// Also a comment", 179, 195},
		{TTSingleWordString, "/+notacomment", 197, 209},
		{TTComment, "/* this is a ml\ncomment that\nshould stop parsing here */", 211, 266},
		{TTMapListBegin, "(", 268, 268},
		{TTMapListEnd, ")", 269, 269},
		{TTEOF, "", 270, 270},
	}

	l := NewTokenizer(input)
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
	input := `/* abort comment on eof`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
		start, end      int
	}{
		{TTComment, "/* abort comment on eof", 0, -1},
	}

	l := NewTokenizer(input)
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
	input := `"abort string on eof`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
		start, end      int
	}{
		{TTString, "abort string on eo", 0, -1},
	}

	l := NewTokenizer(input)
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
	input := ` `

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
		start, end      int
	}{
		{TTEOF, "", 1, 1},
	}

	l := NewTokenizer(input)
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
