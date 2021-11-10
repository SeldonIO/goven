package sql_adaptor

import (
	"errors"
	"fmt"
	"github.com/seldonio/goven/parser"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type (
	// ValidatorFunc takes a field name and validates that it is a legal/correct format.
	ValidatorFunc = func(s string) error
	// ParseValidateFunc takes an Expression from the AST and returns a templated SQL query.
	ParseValidateFunc = func(ex *parser.Expression) (*SqlResponse, error)
)

type SqlResponse struct {
	Raw    string
	Values []string
}

type SqlAdaptor struct {
	// TODO: this feature is not yet implemented.
	// fieldMappings is used to provide ability to map different frontend to backend field names.
	fieldMappings map[string]string
	// defaultFields is the default field matcher, used when a regex isn't matched.
	defaultFields map[string]ParseValidateFunc
	// Non default matchers, these are custom matchers used to extend Goven's functionality.
	matchers map[*regexp.Regexp]ParseValidateFunc
}

func NewSqlAdaptor(fieldMappings map[string]string, defaultFields map[string]ParseValidateFunc, matchers map[*regexp.Regexp]ParseValidateFunc) *SqlAdaptor {
	if fieldMappings == nil {
		fieldMappings = map[string]string{}
	}
	if defaultFields == nil {
		defaultFields = map[string]ParseValidateFunc{}
	}
	if matchers == nil {
		matchers = map[*regexp.Regexp]ParseValidateFunc{}
	}
	sa := SqlAdaptor{
		fieldMappings: fieldMappings,
		defaultFields: defaultFields,
		matchers:      matchers,
	}
	return &sa
}

func (s *SqlAdaptor) Parse(str string) (*SqlResponse, error) {
	newParser := parser.NewParser(str)
	node, err := newParser.Parse()
	if err != nil {
		return nil, errors.New("query could not be parsed")
	}
	return s.parseNodeToSQL(node)
}

func (s *SqlAdaptor) parseNodeToSQL(node parser.Node) (*SqlResponse, error) {
	sq := SqlResponse{}
	if node == nil {
		return &sq, nil
	}
	if node.Type() == parser.EXPRESSION {
		ex, ok := node.(*parser.Expression)
		if !ok {
			return nil, errors.New("failed to parse query correctly")
		}
		// Try and match any custom matchers.
		for k, v := range s.matchers {
			if k.MatchString(ex.Field) {
				return v(ex)
			}
		}
		// If that doesn't happen, then use the relevant default matcher.
		lowerCamelCase := strings.ToLower(strcase.ToCamel(ex.Field))
		if val, ok := s.defaultFields[lowerCamelCase]; ok {
			return val(ex)
		} else {
			// Field is not valid because it must match either a custom regex, or have a validator.
			// If it does neither then we do not expect this field name.
			return nil, errors.New(fmt.Sprintf("field '%s' is not valid", lowerCamelCase))
		}
	}
	op, ok := node.(*parser.Operation)
	if !ok {
		return nil, errors.New("failed to parse query correctly")
	}
	left, err := s.parseNodeToSQL(op.LeftNode)
	if err != nil {
		return nil, err
	}
	// Don't want to have unwanted whitespace if no gate.
	if op.Gate == "" {
		sq = SqlResponse{
			Raw:    fmt.Sprintf("(%s)", left.Raw),
			Values: left.Values,
		}
		return &sq, nil
	}
	right, err := s.parseNodeToSQL(op.RightNode)
	if err != nil {
		return nil, err
	}
	sq = SqlResponse{
		Raw:    fmt.Sprintf("(%s %s %s)", left.Raw, op.Gate, right.Raw),
		Values: append(left.Values, right.Values...),
	}
	return &sq, nil
}
