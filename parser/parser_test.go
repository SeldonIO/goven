package parser

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestBasicParser(t *testing.T) {
	g := NewGomegaWithT(t)
	t.Run("parse expression succeeds", func(t *testing.T) {
		s := "name=max"
		parser := NewParser(s)
		expr, err := parser.parseExpression()
		g.Expect(err).To(BeNil())
		g.Expect(expr).To(Equal(&Expression{
			Field:      "name",
			Comparator: "=",
			Value:      "max",
		}))
	})
	t.Run("parse expression succeeds with whitespace", func(t *testing.T) {
		s := "name=  max "
		parser := NewParser(s)
		expr, err := parser.parseExpression()
		g.Expect(err).To(BeNil())
		g.Expect(expr).To(Equal(&Expression{
			Field:      "name",
			Comparator: "=",
			Value:      "max",
		}))
	})
	t.Run("parse expression fails invalid", func(t *testing.T) {
		s := "name==dog"
		parser := NewParser(s)
		_, err := parser.parseExpression()
		g.Expect(err).ToNot(BeNil())
	})
	t.Run("parse operations badly formatted return errors", func(t *testing.T) {
		tests := []string{
			"name=max AND AND artifact=wow",
			"name=max artifact=wow",
			")(name = max)",
		}
		for _, test := range tests {
			parser := NewParser(test)
			_, err := parser.parseOperation()
			g.Expect(err).ToNot(BeNil(), fmt.Sprintf("failed case: `%s`", test))
		}
	})
	t.Run("parse operations correctly formatted succeeds", func(t *testing.T) {
		test := "name=max AND artifact%art1"
		parser := NewParser(test)
		expected := &Operation{
			LeftNode: &Expression{
				Field:      "name",
				Comparator: "=",
				Value:      "max",
			},
			Gate: "AND",
			RightNode: &Expression{
				Field:      "artifact",
				Comparator: "%",
				Value:      "art1",
			},
		}
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(expected))
	})
	t.Run("parse operations correctly formatted succeeds", func(t *testing.T) {
		test := "(name=max AND artifact=art1) OR metric > 0.98"
		parser := NewParser(test)
		firstExpression := &Operation{
			LeftNode: &Expression{
				Field:      "name",
				Comparator: "=",
				Value:      "max",
			},
			Gate: "AND",
			RightNode: &Expression{
				Field:      "artifact",
				Comparator: "=",
				Value:      "art1",
			},
		}
		secondExpression := &Operation{
			LeftNode: firstExpression,
			Gate:     "OR",
			RightNode: &Expression{
				Field:      "metric",
				Comparator: ">",
				Value:      "0.98",
			},
		}
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(secondExpression))
	})
	t.Run("parse operation when just expression succeeds", func(t *testing.T) {
		test := "name=max"
		parser := NewParser(test)
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(&Expression{
			Field:      "name",
			Comparator: "=",
			Value:      "max",
		}))
	})
	t.Run("parse operation when just bracketed expression succeeds", func(t *testing.T) {
		test := "(name=max)"
		parser := NewParser(test)
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(&Expression{
			Field:      "name",
			Comparator: "=",
			Value:      "max",
		}))
	})
	t.Run("parse operation with metrics/tags format", func(t *testing.T) {
		test := "(metrics[metric-name_1]>= 0.98)"
		parser := NewParser(test)
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(&Expression{
			Field:      "metrics[metric-name_1]",
			Comparator: ">=",
			Value:      "0.98",
		}))
	})
	t.Run("parse operation camelCase", func(t *testing.T) {
		test := "TaskType=classification"
		parser := NewParser(test)
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(&Expression{
			Field:      "TaskType",
			Comparator: "=",
			Value:      "classification",
		}))
	})
	t.Run("parse operation quoted string", func(t *testing.T) {
		test := "(Name=\"Iris Classifier\")"
		parser := NewParser(test)
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(&Expression{
			Field:      "Name",
			Comparator: "=",
			Value:      "Iris Classifier",
		}))
	})
	t.Run("parse empty quoted string", func(t *testing.T) {
		test := "(name=\"\" AND artifact=art1) OR metric > 0.98"
		parser := NewParser(test)
		firstExpression := &Operation{
			LeftNode: &Expression{
				Field:      "name",
				Comparator: "=",
				Value:      "",
			},
			Gate: "AND",
			RightNode: &Expression{
				Field:      "artifact",
				Comparator: "=",
				Value:      "art1",
			},
		}
		secondExpression := &Operation{
			LeftNode: firstExpression,
			Gate:     "OR",
			RightNode: &Expression{
				Field:      "metric",
				Comparator: ">",
				Value:      "0.98",
			},
		}
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(secondExpression))
	})
	t.Run("parsing OR/AND is case insensitive", func(t *testing.T) {
		test := "name=model1 AND version=2.0"
		parser := NewParser(test)
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("AND"))

		test = "name=model1 and version=2.0"
		parser = NewParser(test)
		node, err = parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("AND"))

		test = "name=model1 OR version=2.0"
		parser = NewParser(test)
		node, err = parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("OR"))

		test = "name=model1 or version=2.0"
		parser = NewParser(test)
		node, err = parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("OR"))
	})
	t.Run("parser fails for invalid queries with missing comparator", func(t *testing.T) {
		test := "name"
		parser := NewParser(test)
		_, err := parser.Parse()
		g.Expect(err).ToNot(BeNil())

		test = "name=default OR age"
		parser = NewParser(test)
		_, err = parser.Parse()
		g.Expect(err).ToNot(BeNil())

		test = ""
		parser = NewParser(test)
		_, err = parser.Parse()
		g.Expect(err).ToNot(BeNil())
	})
	t.Run("parser fails for invalid queries with open gate", func(t *testing.T) {
		test := "name=default AND"
		parser := NewParser(test)
		_, err := parser.Parse()
		g.Expect(err).ToNot(BeNil())
	})
}

func FuzzParser(f *testing.F) {
	testcases := []string{">=!=", "name=default OR age", "< <= = != AND OR and or", "1  !=   \"2\"", "(Name=\"Iris Classifier\")"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, s string) {
		parser := NewParser(s)
		node, err := parser.Parse()
		if err == nil {
			_, nodeIsOp := node.(*Operation)
			_, nodeIsExpr := node.(*Expression)
			if !nodeIsOp && !nodeIsExpr {
				t.Errorf("node must be either op or expression")
			}
		}
	})
}
