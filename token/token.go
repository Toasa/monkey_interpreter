package token

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}

const (
    // TokenTypeの種類

    ILLGAL = "ILLGAL"
    EOF = "EOF"

    // ident, literal
    IDENT = "IDENT" // add, foo, x, y,
    STRING = "STRING"
    INT = "INT" // 46

    // operator
    ASSGIN = "="
    PLUS = "+"
    MINUS = "-"
    MUL = "*"
    DIV = "/"
    LT = "<"
    GT = ">"
    BANG = "!"

    // delimiter
    COMMA = ","
    COLON = ":"
    SEMICOLON = ";"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"
    LBRACKET = "["
    RBRACKET = "]"

    EQ = "=="
    NQ = "!="

    // keyword
    FUNCTION = "FUNCTION"
    LET = "LET"
    TRUE = "TRUE"
    FALSE = "FALSE"
    IF = "IF"
    ELSE = "ELSE"
    RETURN = "RETURN"
)

var keywords = map[string]TokenType {
    "let": LET,
    "fn": FUNCTION,
    "true": TRUE,
    "false": FALSE,
    "if": IF,
    "else": ELSE,
    "return": RETURN,
}

func LookupIdent(str string) TokenType {
    keyword, ok := keywords[str]
    if ok {
        return keyword
    }
    return IDENT
}
