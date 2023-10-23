package lexer

import (
	"github.com/dominiwe/jsonfs/better/token"
	"unicode/utf8"
)

type Lexer struct {
	input       []byte
	inputLength int
	position    int
	rune        rune
	eof         bool
}

func New(input []byte) *Lexer {
	lexer := &Lexer{input: input}
	lexer.inputLength = len(input)
	lexer.eof = false
	lexer.consumeRune()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.eatWhitespace()

	switch lexer.rune {
	case '{':
		tok = newToken(token.LBRACE, lexer.rune)
	case '}':
		tok = newToken(token.RBRACE, lexer.rune)
	case '[':
		tok = newToken(token.LBRACK, lexer.rune)
	case ']':
		tok = newToken(token.RBRACK, lexer.rune)
	case ':':
		tok = newToken(token.COLON, lexer.rune)
	case ',':
		tok = newToken(token.COMMA, lexer.rune)
	case '"':
		var legal bool
		tok.Type = token.STRING
		tok.Literal, legal = lexer.consumeString()
		if !legal {
			tok.Type = token.ILLEGAL
		}
	default:
		if lexer.startOfNumber() {
			var legal bool
			tok.Type = token.NUMBER
			tok.Literal, legal = lexer.consumeNumber()
			if !legal {
				tok.Type = token.ILLEGAL
			}
		} else if isOtherLiteral(lexer.rune) {
			var legal bool
			switch lexer.rune {
			case 'n':
				tok.Type = token.NULL
				tok.Literal, legal = lexer.consumeOtherLit("null")
			case 't':
				tok.Type = token.TRUE
				tok.Literal, legal = lexer.consumeOtherLit("true")
			case 'f':
				tok.Type = token.FALSE
				tok.Literal, legal = lexer.consumeOtherLit("false")
			}
			if !legal {
				tok.Type = token.ILLEGAL
			}
		} else if lexer.eof {
			tok.Type = token.EOF
			tok.Literal = ""
		} else {
			tok = newToken(token.ILLEGAL, lexer.rune)
		}
	}

	lexer.consumeRune()
	return tok
}

func (lexer *Lexer) eatWhitespace() {
	for isWs(lexer.rune) {
		lexer.consumeRune()
	}
}

func newToken(tokenType token.Type, r rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(r)}
}

func (lexer *Lexer) consumeRune() {
	decodedRune, size := utf8.DecodeRune(lexer.input[lexer.position:])
	if size == 0 {
		lexer.eof = true
	}
	lexer.position += size
	lexer.rune = decodedRune
}

func (lexer *Lexer) peekRune() (rune, int) {
	return utf8.DecodeRune(lexer.input[lexer.position:])
}

func isWs(r rune) bool {
	return r == '\u0020' || // space
		r == '\u000A' || // line feed
		r == '\u000D' || // carriage return
		r == '\u0009' // horizontal tab
}

func isSign(r rune) bool {
	return r == '+' || r == '-'
}

func isMinus(r rune) bool {
	return r == '-'
}

func isExponent(r rune) bool {
	return r == 'e' || r == 'E'
}

func isOneNine(r rune) bool {
	return '1' <= r && r <= '9'
}

func isDigit(r rune) bool {
	return isOneNine(r) || r == '0'
}

func isHex(r rune) bool {
	return isDigit(r) || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F'
}

func isChar(r rune) bool {
	return r != '"' && r != '\\' && r > '\u001F'
}

func isFraction(r rune) bool {
	return r == '.'
}

func isOtherLiteral(r rune) bool {
	return r == 'n' || r == 't' || r == 'f'
}

func isEscapedChar(r rune) bool {
	return r == '"' ||
		r == '\\' ||
		r == '/' ||
		r == 'b' ||
		r == 'f' ||
		r == 'n' ||
		r == 'r' ||
		r == 't'
}

func (lexer *Lexer) consumeString() (string, bool) {
	position := lexer.position
	legal := true
	offset := 1
out:
	for {
		lexer.consumeRune()
		if lexer.eof {
			legal = false
			break
		} else if lexer.rune == '"' {
			break
		} else if lexer.rune == '\\' {
			lexer.consumeRune()
			if lexer.rune == 'u' {
				for digit := 1; digit <= 4; digit++ {
					lexer.consumeRune()
					if !isHex(lexer.rune) {
						legal = false
						break out
					}
				}
			} else if isEscapedChar(lexer.rune) {
				lexer.consumeRune()
			} else {
				legal = false
				break
			}
		} else if !isChar(lexer.rune) {
			legal = false
			break
		}
	}
	if !legal {
		offset = 0
	}
	return string(lexer.input[position:(lexer.position - offset)]), legal
}

func (lexer *Lexer) consumeFraction() bool {
	r, _ := lexer.peekRune()
	if !isDigit(r) {
		return false
	} else {
		lexer.consumeRune()
	}
	for {
		r, _ = lexer.peekRune()
		if isExponent(r) {
			lexer.consumeRune()
			return lexer.consumeExponent()
		} else if !isDigit(r) {
			return true
		} else {
			lexer.consumeRune()
		}
	}
}

func (lexer *Lexer) consumeExponent() bool {
	r, _ := lexer.peekRune()
	if !isSign(r) {
		return false
	} else {
		lexer.consumeRune()
	}
	r, _ = lexer.peekRune()
	if !isDigit(r) {
		return false
	} else {
		lexer.consumeRune()
	}
	for {
		r, _ = lexer.peekRune()
		if !isDigit(r) {
			return true
		} else {
			lexer.consumeRune()
		}
	}
}

func (lexer *Lexer) consumeNumber() (string, bool) {
	position := lexer.position - 1
	legal := true
	if isMinus(lexer.rune) {
		lexer.consumeRune()
	}
	for {
		r, _ := lexer.peekRune()
		if isFraction(r) {
			lexer.consumeRune()
			legal = lexer.consumeFraction()
			break
		} else if isExponent(r) {
			lexer.consumeRune()
			legal = lexer.consumeExponent()
			break
		} else if !isDigit(r) {
			break
		} else {
			lexer.consumeRune()
		}
	}
	return string(lexer.input[position:lexer.position]), legal
}

func (lexer *Lexer) consumeOtherLit(literal string) (string, bool) {
	position := lexer.position - 1
	legal := true
	for _, char := range literal[1:] {
		lexer.consumeRune()
		if lexer.eof || lexer.rune != char {
			legal = false
			break
		}
	}
	return string(lexer.input[position:lexer.position]), legal
}

func (lexer *Lexer) startOfNumber() bool {
	if isMinus(lexer.rune) {
		r, _ := lexer.peekRune()
		return isDigit(r)
	} else {
		return isDigit(lexer.rune)
	}
}
