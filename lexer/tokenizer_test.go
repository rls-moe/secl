package lexer

import "testing"

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

	tests := []struct{
		expectedType TokenType
		expectedLiteral string
		start,end int
	}{
		{TT_mapListBegin, "(", 0, 0},
		{TT_mod_execMap, "!", 2, 2},
		{TT_mapListBegin, "(", 3, 3},
		{TT_function, "nop", 4, 6},
		{TT_mapListEnd, ")", 7, 7},
		{TT_singleWordString, "let", 9, 11},
		{TT_mod_mapKey, ":", 12, 12},
		{TT_number, "909", 14, 16},
		{TT_singleWordString, "other", 18, 22},
		{TT_mod_mapKey, ":", 23, 23},
		{TT_mod_trim, "@", 24, 24},
		{TT_string, "test", 25, 30},
		{TT_singleWordString, "k2", 32, 33},
		{TT_mod_mapKey, ":", 34, 34},
		{TT_string, "test\n\thello🤣🤣\\\"ttt", 36, 55},
		{TT_singleWordString, "k3?", 57, 59},
		{TT_mod_mapKey, ":", 60, 60},
		{TT_singleWordString, "k4!", 62, 64},
		{TT_singleWordString, "decimal", 66, 72},
		{TT_mod_mapKey, ":", 73, 73},
		{TT_number, "0.111e-19*10^10", 75, 89},
		{TT_singleWordString, "let2", 91, 94},
		{TT_mod_mapKey, ":", 95, 95},
		{TT_empty, "empty", 97, 101},
		{TT_singleWordString, "let3", 103, 106},
		{TT_mod_mapKey, ":", 107, 107},
		{tt_nil, "nil",109,111},
		{TT_randstr, "randstr128", 113, 122},
		{TT_function, "decb64", 124, 129},
		{TT_bool, "true", 131, 134},
		{TT_bool, "off", 136, 138},
		{TT_mapListEnd, ")", 140, 140},
		{TT_comment, "# This is a comment", 142, 160},
		{TT_comment, "; Also a comment", 162, 177},
		{TT_comment, "// Also a comment", 179, 195},
		{TT_singleWordString, "/+notacomment", 197, 209},
		{TT_comment, "/* this is a ml\ncomment that\nshould stop parsing here */", 211, 266},
		{TT_mapListBegin, "(", 268, 268},
		{TT_mapListEnd, ")", 269, 269},
		{TT_eof, "", 270, 270},
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
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i , tests[i].start, tests[i].end, tok.Start, tok.End)
		}

		//t.Logf("Token: %s", tok.Type)
	}
}

func TestTokenizer_NextToken2(t *testing.T) {
	input := `/* abort comment on eof`

	tests := []struct{
		expectedType TokenType
		expectedLiteral string
		start,end int
	}{
		{TT_comment, "/* abort comment on eof", 0, -1},
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
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i , tests[i].start, tests[i].end, tok.Start, tok.End)
		}
	}
}

func TestTokenizer_NextToken3(t *testing.T) {
	input := `"abort string on eof`

	tests := []struct{
		expectedType TokenType
		expectedLiteral string
		start,end int
	}{
		{TT_string, "abort string on eo", 0, -1},
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
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i , tests[i].start, tests[i].end, tok.Start, tok.End)
		}
	}
}


func TestTokenizer_NextToken4(t *testing.T) {
	input := ` `

	tests := []struct{
		expectedType TokenType
		expectedLiteral string
		start,end int
	}{
		{TT_eof, "", 1, 1},
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
			t.Fatalf("test %d: position wrong, expected %d-%d got %d-%d", i , tests[i].start, tests[i].end, tok.Start, tok.End)
		}
	}
}