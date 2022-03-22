package sql_adaptor

import (
	"fmt"
	"strconv"

	"github.com/seldonio/goven/parser"
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
	if ex.Comparator == parser.TokenLookup[parser.PERCENT] {
		fmtValue := fmt.Sprintf("%%%s%%", ex.Value)
		sq := SqlResponse{
			Raw:    fmt.Sprintf("%s LIKE ?", ex.Field),
			Values: []string{fmtValue},
		}
		return &sq
	}
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
		return fmt.Errorf("value '%s' is not an integer", s)
	}
	return nil
}

func NumericValidator(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("value '%s' is not numeric", s)
	}
	return nil
}
