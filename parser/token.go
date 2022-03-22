package parser

// Token represents a lexical token.
type Token int

// TokenInfo stores relevant information about the token during scanning.
type TokenInfo struct {
	Token   Token
	Literal string
}

// TokenLookup is a map, useful for printing readable names of the tokens.
var TokenLookup = map[Token]string{
	OTHER:               "OTHER",
	EOF:                 "EOF",
	WS:                  "WS",
	STRING:              "STRING",
	EQUAL:               "EQUAL",
	GREATER_THAN:        "GREATER THAN",
	GREATHER_THAN_EQUAL: "GREATER THAN OR EQUAL",
	LESS_THAN:           "LESS THAN",
	LESS_THAN_EQUAL:     "LESS THAN OR EQUAL",
	NOT_EQUAL:           "NOT EQUAL",
	AND:                 "AND",
	OR:                  "OR",
	OPEN_BRACKET:        "(",
	CLOSED_BRACKET:      ")",
	PERCENT:             "%",
}

// String prints a human readable string name for a given token.
func (t Token) String() (print string) {
	return TokenLookup[t]
}

// Declare the tokens here.
const (
	// Special tokens
	// Iota simply starts and integer count
	OTHER Token = iota
	EOF
	WS

	// Main literals
	STRING

	// Brackets
	OPEN_BRACKET
	CLOSED_BRACKET

	// Special characters
	GREATER_THAN
	GREATHER_THAN_EQUAL
	LESS_THAN
	LESS_THAN_EQUAL
	EQUAL
	NOT_EQUAL
	PERCENT

	// Keywords
	AND
	OR
)
