package sql_adaptor_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/seldonio/goven/sql_adaptor"

	. "github.com/onsi/gomega"
)

type TestCase struct {
	test           string
	expectedRaw    string
	expectedValues []string
}

// Typical gorm database struct - a User
type ExampleDBStruct struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func TestSqlAdaptor(t *testing.T) {
	g := NewGomegaWithT(t)
	t.Run("test sql adaptor success", func(t *testing.T) {
		testCases := []TestCase{
			{
				test:           "(name=max AND email=bob-dylan@aol.com) OR age > 1",
				expectedRaw:    "((name=? AND email=?) OR age>?)",
				expectedValues: []string{"max", "bob-dylan@aol.com", "1"},
			},
			// Test for an empty quoted string.
			{
				test:           "(name=\"\" AND email=bob-dylan@aol.com) OR age > 1",
				expectedRaw:    "((name=? AND email=?) OR age>?)",
				expectedValues: []string{"", "bob-dylan@aol.com", "1"},
			},
			{
				test:           "(name%max AND email=bob-dylan@aol.com) OR age > 1",
				expectedRaw:    "((name LIKE ? AND email=?) OR age>?)",
				expectedValues: []string{"%max%", "bob-dylan@aol.com", "1"},
			},
		}
		for _, testCase := range testCases {
			sa := sql_adaptor.NewDefaultAdaptorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
			response, err := sa.Parse(testCase.test)
			g.Expect(err).To(BeNil(), fmt.Sprintf("failed case: %s", testCase.test))
			g.Expect(response.Raw).To(Equal(testCase.expectedRaw), fmt.Sprintf("failed case raw: %s", testCase.test))
			g.Expect(response.Values).To(Equal(testCase.expectedValues), fmt.Sprintf("failed case values: %s", testCase.test))
		}
	})
	t.Run("test sql adaptor failure", func(t *testing.T) {
		testCases := []TestCase{
			{
				test: "(name=max AND invalidField=wow) OR age > 1",
			},
			{
				test: "id = wow",
			},
			{
				test: "age = wow",
			},
			{
				test: "name = default AND",
			},
			{
				test: "name",
			},
			{
				test: "name = default AND age",
			},
		}
		for _, testCase := range testCases {
			sa := sql_adaptor.NewDefaultAdaptorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
			_, err := sa.Parse(testCase.test)
			g.Expect(err).ToNot(BeNil(), fmt.Sprintf("failed case: %s", testCase.test))
		}
	})
	t.Run("test FieldParseValidatorFromStruct", func(t *testing.T) {
		defaultFields := sql_adaptor.FieldParseValidatorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
		_, ok := defaultFields["membernumber"]
		g.Expect(ok).To(Equal(true))
	})
}

func FuzzSqlAdaptor(f *testing.F) {
	testcases := []string{"(name=max AND invalidField=wow) OR age > 1", "(name=max AND email=bob-dylan@aol.com) OR age > 1", "id = wow", "(name%max AND email=\"bob-dylan@aol.com\") OR age > 1"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, s string) {
		sa := sql_adaptor.NewDefaultAdaptorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
		response, err := sa.Parse(s)
		if err != nil && response != nil {
			t.Errorf("expected nil response when err is not nil")
		}
	})
}
