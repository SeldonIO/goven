package parser

import (
	"testing"

	. "github.com/onsi/gomega"
)

func lexerHelper(lex *Lexer) ([]Token, []string) {
	var tokens []Token
	var literals []string
	for {
		tok := lex.Scan()
		tokens = append(tokens, tok.Token)
		literals = append(literals, tok.Literal)
		if tok.Token == EOF {
			return tokens, literals
		}
	}
}

func TestLexer(t *testing.T) {
	g := NewGomegaWithT(t)
	t.Run("scan into tokens succeeds", func(t *testing.T) {
		s := "name=model1"
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{STRING, EQUAL, STRING, EOF}))
		g.Expect(literals).To(Equal([]string{"name", "=", "model1", ""}))
	})
	t.Run("scan into tokens succeeds percent", func(t *testing.T) {
		s := "name%model1"
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{STRING, PERCENT, STRING, EOF}))
		g.Expect(literals).To(Equal([]string{"name", "%", "model1", ""}))
	})
	t.Run("scan into tokens succeeds for quoted string", func(t *testing.T) {
		s := "name=\"Iris Classifier\""
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{STRING, EQUAL, STRING, EOF}))
		g.Expect(literals).To(Equal([]string{"name", "=", "Iris Classifier", ""}))
	})
	t.Run("scan into tokens succeeds for empty quoted string", func(t *testing.T) {
		s := "name=\"\""
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{STRING, EQUAL, STRING, EOF}))
		g.Expect(literals).To(Equal([]string{"name", "=", "", ""}))
	})
	t.Run("scan into tokens with whitespace succeeds", func(t *testing.T) {
		s := "name =   \n   model1"
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{STRING, WS, EQUAL, WS, STRING, EOF}))
		g.Expect(literals).To(Equal([]string{"name", "", "=", "", "model1", ""}))
	})
	t.Run("scan into tokens all token types", func(t *testing.T) {
		s := "string ( ) > >= < <= = != AND OR and or"
		lexer := NewLexerFromString(s)
		tokens, literals := lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{STRING, WS, OPEN_BRACKET, WS, CLOSED_BRACKET, WS, GREATER_THAN, WS,
			GREATHER_THAN_EQUAL, WS, LESS_THAN, WS, LESS_THAN_EQUAL, WS, EQUAL, WS, NOT_EQUAL, WS, AND, WS, OR, WS, AND, WS, OR, EOF}))
		var literalsNoWhitespace []string
		for _, val := range literals {
			if val != "" {
				literalsNoWhitespace = append(literalsNoWhitespace, val)
			}
		}
		g.Expect(literalsNoWhitespace).To(Equal([]string{"string", "(", ")", ">", ">=", "<", "<=", "=", "!=", "AND", "OR", "AND", "OR"}))
	})
	t.Run("scan tokens is greedy", func(t *testing.T) {
		s := "<=="
		lexer := NewLexerFromString(s)
		tokens, _ := lexerHelper(lexer)
		// Here we expect the string to be consumed as less than equal, equal. Not less than, equal, equal.
		g.Expect(tokens).To(Equal([]Token{LESS_THAN_EQUAL, EQUAL, EOF}))
		s = ">=="
		lexer = NewLexerFromString(s)
		tokens, _ = lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{GREATHER_THAN_EQUAL, EQUAL, EOF}))
		s = "!=="
		lexer = NewLexerFromString(s)
		tokens, _ = lexerHelper(lexer)
		g.Expect(tokens).To(Equal([]Token{NOT_EQUAL, EQUAL, EOF}))
	})
}

func FuzzLexer(f *testing.F) {
	testcases := []string{">=!=", "string ( ) > >=", "< <= = != AND OR and or", "1  !=   \"2\""}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, s string) {
		lexer := NewLexerFromString(s)
		tokens, _ := lexerHelper(lexer)
		// This is really testing for panics only.
		for _, token := range tokens {
			if _, ok := TokenLookup[token]; !ok {
				t.Errorf("unexpected token %v", token)
			}
		}
	})
}
