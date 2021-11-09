package parser

import (
	"errors"
)

// TokenStack is used as the buffer for the Parser.
type TokenStack struct {
	stack []TokenInfo
}

func (s *TokenStack) Push(v TokenInfo) {
	s.stack = append(s.stack, v)
}

func (s *TokenStack) Pop() (TokenInfo, error) {
	if len(s.stack) == 0 {
		return TokenInfo{}, errors.New("stack is empty")
	}
	l := len(s.stack)
	token := s.stack[l-1]
	s.stack = s.stack[:l-1]
	return token, nil
}

func (s *TokenStack) Len() int {
	return len(s.stack)
}
