package token

type Type string

type Token struct {
	T       Type
	Literal string
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"null":   NULL,
	"import": IMPORT,
	"for":    FOR,
	"in":     IN,
	"break":  BREAK,
}

func New(tokType Type, literal string) Token {
	return Token{T: tokType, Literal: literal}
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

const (
	// Special, non-visible tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers, literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN = "="
	// Arithmetic
	PLUS      = "+"
	MINUS     = "-"
	TIMES     = "*"
	DIVIDE    = "/"
	REMAINDER = "%"
	// Boolean operators
	AND = "&&"
	OR  = "||"
	// Comparison
	GREATER_THAN = ">"
	GEQ          = ">="
	LESS_THAN    = "<"
	LEQ          = "<="
	BANG         = "!"
	EQUAL        = "=="
	NOT_EQUAL    = "!="
	// Increment/Decrement
	INCR     = "+="
	INCR_ONE = "++"
	DECR     = "-="
	DECR_ONE = "--"
	// Other special operators that have some use
	TIMES_EQ = "*="
	DIV_EQ   = "/="
	POWER    = "**"
	REM_EQ   = "%="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// Accessors
	PERIOD = "."

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACK = "["
	RBRACK = "]"

	// Reserved words
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
	NULL     = "NULL"
	IMPORT   = "IMPORT"
	FOR      = "FOR"
	IN       = "IN"
	BREAK    = "BREAK"
)
