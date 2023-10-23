package lexer

import (
	"github.com/dominiwe/jsonfs/better/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := []byte(`{"fooŧ": "bar\t \ndwa", "baz\uAAD2": [1.021, -2, 3e+2]}{,.l**"d\i"}`)
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LBRACE, "{"},
		{token.STRING, "fooŧ"},
		{token.COLON, ":"},
		{token.STRING, "bar\\t \\ndwa"},
		{token.COMMA, ","},
		{token.STRING, "baz\\uAAD2"},
		{token.COLON, ":"},
		{token.LBRACK, "["},
		{token.NUMBER, "1.021"},
		{token.COMMA, ","},
		{token.NUMBER, "-2"},
		{token.COMMA, ","},
		{token.NUMBER, "3e+2"},
		{token.RBRACK, "]"},
		{token.RBRACE, "}"},
		{token.LBRACE, "{"},
		{token.COMMA, ","},
		{token.ILLEGAL, "."},
		{token.ILLEGAL, "l"},
		{token.ILLEGAL, "*"},
		{token.ILLEGAL, "*"},
		{token.ILLEGAL, "d\\i"},
		{token.ILLEGAL, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		//fmt.Printf("t: \"%s\"\n\tv: \"%s\"\n", tok.Type, tok.Literal)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken2(t *testing.T) {
	input := []byte(`{"foo":{"foo": 123.123e-3132, "bar":0.02E+32}, "bar":null, "baz": true, "baw": false, "baw2": {}, "baw3": []}`)
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.NUMBER, "123.123e-3132"},
		{token.COMMA, ","},
		{token.STRING, "bar"},
		{token.COLON, ":"},
		{token.NUMBER, "0.02E+32"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.STRING, "bar"},
		{token.COLON, ":"},
		{token.NULL, "null"},
		{token.COMMA, ","},
		{token.STRING, "baz"},
		{token.COLON, ":"},
		{token.TRUE, "true"},
		{token.COMMA, ","},
		{token.STRING, "baw"},
		{token.COLON, ":"},
		{token.FALSE, "false"},
		{token.COMMA, ","},
		{token.STRING, "baw2"},
		{token.COLON, ":"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.STRING, "baw3"},
		{token.COLON, ":"},
		{token.LBRACK, "["},
		{token.RBRACK, "]"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		//fmt.Printf("t: \"%s\"\n\tv: \"%s\"\n", tok.Type, tok.Literal)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken3(t *testing.T) {
	inputs := [][]byte{
		[]byte(`{`),
		[]byte(`{"f`),
		[]byte(``),
		[]byte(`{"foo":nu`),
		[]byte(`{"foo":tr`),
		[]byte(`{"foo":fal`),
		[]byte(`{"foo":1.231`),
		[]byte(`{"foo":1.231e`),
		[]byte(`{"foo":1.231e-`),
		[]byte(`{"foo":1.231e-1`),
		[]byte(`{"foo":[`),
		[]byte(`{"foo":"dwadaw`),
	}
	tests := [][]struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{
			{token.LBRACE, "{"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.ILLEGAL, "f"},
			{token.EOF, ""},
		},
		{
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.ILLEGAL, "nu"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.ILLEGAL, "tr"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.ILLEGAL, "fal"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.NUMBER, "1.231"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.ILLEGAL, "1.231e"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.ILLEGAL, "1.231e-"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.NUMBER, "1.231e-1"},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.LBRACK, "["},
			{token.EOF, ""},
		},
		{
			{token.LBRACE, "{"},
			{token.STRING, "foo"},
			{token.COLON, ":"},
			{token.ILLEGAL, "dwadaw"},
			{token.EOF, ""},
		},
	}
	for i, input := range inputs {
		l := New(input)
		test := tests[i]
		for j, tt := range test {
			tok := l.NextToken()
			//fmt.Printf("t: \"%s\"\n\tv: \"%s\"\n", tok.Type, tok.Literal)

			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					j, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
					j, tt.expectedLiteral, tok.Literal)
			}
		}
	}
}
