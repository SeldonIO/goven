package parser

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTokenStack(t *testing.T) {
	g := NewGomegaWithT(t)
	t.Run("test token stack", func(t *testing.T) {
		stack := TokenStack{
			stack: []TokenInfo{},
		}
		stack.Push(TokenInfo{
			Token:   EQUAL,
			Literal: "=",
		})
		length := stack.Len()
		g.Expect(length).To(Equal(1))
		tok, err := stack.Pop()
		g.Expect(err).To(BeNil())
		g.Expect(tok).To(Equal(TokenInfo{
			Token:   EQUAL,
			Literal: "=",
		}))
	})
}
