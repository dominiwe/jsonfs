package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	LBRACE  = "{"
	RBRACE  = "}"
	LBRACK  = "["
	RBRACK  = "]"
	COLON   = ":"
	COMMA   = ","
	STRING  = "STRING"
	NUMBER  = "NUMBER"
	TRUE    = "TRUE"
	FALSE   = "FALSE"
	NULL    = "NULL"
)
