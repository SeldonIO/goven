package parser

import "fmt"

const (
	OPERATION  = "operation"
	EXPRESSION = "expression"
)

// Node represents a node in the AST after the expression is parsed.
type Node interface {
	Type() string
	String() string
}

// Expression represents something like x=y or x>=y
type Expression struct {
	Field      string
	Comparator string
	Value      string
}

// Operation represents a Node (Operation or Expression) compared with another Node using either `AND` or `OR`.
type Operation struct {
	LeftNode  Node
	Gate      string
	RightNode Node
}

func (e Expression) Type() string { return EXPRESSION }
func (o Operation) Type() string  { return OPERATION }

func (e Expression) String() string {
	return fmt.Sprintf("%v %v %v", e.Field, e.Comparator, e.Value)
}

func (o Operation) String() string {
	if o.Gate == "" {
		return fmt.Sprintf("(%v)", o.LeftNode)
	}
	return fmt.Sprintf("(%v %v %v)", o.LeftNode, o.Gate, o.RightNode)
}
