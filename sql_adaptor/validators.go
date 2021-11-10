package sql_adaptor

import (
	"errors"
	"fmt"
	"github.com/seldonio/goven/parser"
	"strconv"
)

// DefaultMatcherWithValidator wraps the default matcher with validation on the value.
func DefaultMatcherWithValidator(validate ValidatorFunc) ParseValidateFunc {
	return func(ex *parser.Expression) (*SqlResponse, error) {
		err := validate(ex.Value)
		if err != nil {
			return nil, err
		}
		return DefaultMatcher(ex), nil
	}
}

// DefaultMatcher takes an expression and spits out the default SqlResponse.
func DefaultMatcher(ex *parser.Expression) *SqlResponse {
	sq := SqlResponse{
		Raw:    fmt.Sprintf("%s%s?", ex.Field, ex.Comparator),
		Values: []string{ex.Value},
	}
	return &sq
}

func NullValidator(_ string) error {
	return nil
}

func IntegerValidator(s string) error {
	_, err := strconv.Atoi(s)
	if err != nil {
		return errors.New(fmt.Sprintf("value '%s' is not an integer", s))
	}
	return nil
}

func NumericValidator(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("value '%s' is not numeric", s))
	}
	return nil
}
