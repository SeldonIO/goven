package parser

import (
	"errors"
)

// TokenStack is used as the buffer for the Parser.
type TokenStack struct {
	stack []TokenInfo
}

// Push pushes a token to the TokenStack.
func (s *TokenStack) Push(v TokenInfo) {
	s.stack = append(s.stack, v)
}

// Pop removes and returns a token from the TokenStack.
func (s *TokenStack) Pop() (TokenInfo, error) {
	if len(s.stack) == 0 {
		return TokenInfo{}, errors.New("stack is empty")
	}
	l := len(s.stack)
	token := s.stack[l-1]
	s.stack = s.stack[:l-1]
	return token, nil
}

// Len returns the current length of the TokenStack.
func (s *TokenStack) Len() int {
	return len(s.stack)
}
