package parser

import (
	"errors"
	"fmt"
	"strings"
)

// Parser represents a parser, including a scanner and the underlying raw input.
// It also contains a small buffer to allow for two unscans.
type Parser struct {
	s   *Lexer
	raw string
	buf TokenStack
}

// NewParser returns a new instance of Parser.
func NewParser(s string) *Parser {
	return &Parser{s: NewLexer(strings.NewReader(s)), raw: s}
}

// Parse takes the raw string and returns the root node of the AST.
func (p *Parser) Parse() (Node, error) {
	operation, err := p.parseOperation()
	if err != nil {
		return nil, err
	}
	// Try to peel like an onion.
	gate := operation.(*Operation).Gate
	if gate != "" && operation.(*Operation).RightNode == nil {
		return nil, errors.New("found open gate")
	}
	for gate == "" {
		operation = operation.(*Operation).LeftNode
		if operation == nil {
			return nil, errors.New("got nil operation")
		}
		if operation.Type() == EXPRESSION {
			break
		}
		gate = operation.(*Operation).Gate
	}
	return operation, nil
}

func (p *Parser) parseOperation() (Node, error) {
	op := &Operation{
		LeftNode:  nil,
		Gate:      "",
		RightNode: nil,
	}
	tok, lit := p.scanIgnoreWhitespace()
	for tok != EOF {
		switch {
		// If we hit an open bracket then we parse the operation contained in the brackets.
		case tok == OPEN_BRACKET:
			node, err := p.parseOperation()
			if err != nil {
				return nil, err
			}
			// Assign the operation to left node if we haven't already.
			if op.LeftNode == nil {
				op.LeftNode = node
				break
			}
			if op.Gate == "" {
				return nil, errors.New("shouldn't find operation before Gate if left node already exists")
			}
			// Assign to right otherwise.
			if op.RightNode == nil {
				op.RightNode = node
				tempOp := &Operation{
					LeftNode:  op,
					Gate:      "",
					RightNode: nil,
				}
				op = tempOp
				break
			}
		case tok == STRING:
			if (op.LeftNode != nil && op.Gate == "") || (op.LeftNode != nil && op.RightNode != nil) {
				return nil, errors.New("didn't expect an expression here")
			}
			p.unscan(TokenInfo{
				Token:   tok,
				Literal: lit,
			})
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if op.LeftNode == nil {
				op.LeftNode = expr
				break
			}
			if op.RightNode == nil {
				op.RightNode = expr
				tempOp := &Operation{
					LeftNode:  op,
					Gate:      "",
					RightNode: nil,
				}
				op = tempOp
				break
			}
			return nil, errors.New("parsed expression, left and right node shouldn't have both been populated")
		case tok == CLOSED_BRACKET:
			if op.LeftNode == nil {
				return nil, errors.New("can't close a bracket when we have parsed nothing for the left node")
			}
			return op, nil
		case isTokenGate(tok):
			if op.Gate != "" {
				return nil, errors.New("already found a Gate")
			}
			op.Gate = lit
		default:
			return nil, fmt.Errorf("unexpected token %v", tok)
		}
		tok, lit = p.scanIgnoreWhitespace()
	}
	return op, nil
}

func (p *Parser) parseExpression() (Node, error) {
	exp := &Expression{
		Field:      "",
		Comparator: "",
		Value:      "",
	}

	isValueEmpty := false
	// This code relies on being Field -> Comparator -> Value in order.
	tok, lit := p.scan()
	for tok != EOF {
		switch {
		// Open bracket means we found an operation that needs parsing.
		case tok == OPEN_BRACKET:
			return p.parseOperation()
		// Ignore whitespace unless we have completed the expression.
		case tok == WS:
			if exp.Field != "" && exp.Comparator != "" && (exp.Value != "" || isValueEmpty) {
				// Got to the end of the expression so quit
				return exp, nil
			}
		// Looking for the Field name.
		case exp.Field == "":
			if tok != STRING {
				return nil, fmt.Errorf("expected Field, got %v", tok)
			}
			exp.Field = lit
		// Looking for the Comparator.
		case exp.Comparator == "":
			if !isTokenComparator(tok) {
				return nil, fmt.Errorf("expected Comparator, got %v", tok)
			}
			exp.Comparator = lit
		// Looking for the Value
		case exp.Value == "":
			if tok != STRING {
				// If we didn't have an empty string in the value field - return an error.
				if isValueEmpty {
					break
				} else {
					return nil, fmt.Errorf("expected Value, got %v", tok)
				}
			}
			if lit == "" {
				isValueEmpty = true
			}
			exp.Value = lit
		}
		tok, lit = p.scan()
	}
	if exp.Field != "" && exp.Comparator == "" {
		return nil, errors.New("found no comparator when expected")
	}
	return exp, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.Len() != 0 {
		// Can ignore the error since it's not empty.
		tokenInf, _ := p.buf.Pop()
		return tokenInf.Token, tokenInf.Literal
	}

	// Otherwise read the next token from the scanner.
	tokenInf := p.s.Scan()
	tok, lit = tokenInf.Token, tokenInf.Literal
	return tok, lit
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return tok, lit
}

// unscan pushes the previously read tokens back onto the buffer.
func (p *Parser) unscan(tok TokenInfo) {
	p.buf.Push(tok)
}
