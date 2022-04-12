// Package example provides example use cases of goven with a data model.
package example

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/seldonio/goven/parser"

	"github.com/seldonio/goven/sql_adaptor"
	"gorm.io/gorm"
)

const (
	KeyValueRegex = `(.+)\[(.+)\]`
)

type Model struct {
	gorm.Model
	Name      string
	Version   string
	CreatedAt time.Time
	Tags      []Tag `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Tag struct {
	gorm.Model
	Key     string
	Value   string
	ModelID uint
}

type ModelDAO struct {
	db           *gorm.DB
	queryAdaptor *sql_adaptor.SqlAdaptor
}

func NewModelDAO(db *gorm.DB) (*ModelDAO, error) {
	adaptor, err := CreateModelAdaptor()
	if err != nil {
		return nil, err
	}
	return &ModelDAO{
		db:           db,
		queryAdaptor: adaptor,
	}, nil
}

func (u *ModelDAO) CreateModel(model *Model) error {
	ctx := context.Background()
	tx := u.db.Begin().WithContext(ctx)
	err := tx.Create(model).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (u *ModelDAO) MakeQuery(q string) ([]Model, error) {
	var models []Model
	ctx := context.Background()
	query := u.db.WithContext(ctx)
	queryResp, err := u.queryAdaptor.Parse(q)
	if err != nil {
		return nil, err
	}
	query = query.Preload("Tags").Model(Model{}).Where(queryResp.Raw, sql_adaptor.StringSliceToInterfaceSlice(queryResp.Values)...)
	err = query.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func CreateModelAdaptor() (*sql_adaptor.SqlAdaptor, error) {
	matchers := map[*regexp.Regexp]sql_adaptor.ParseValidateFunc{}
	fieldMappings := map[string]string{}

	// Custom matcher initialised here.
	reg, err := regexp.Compile(KeyValueRegex)
	if err != nil {
		return nil, err
	}
	matchers[reg] = keyValueMatcher

	reflection := reflect.ValueOf(&Model{})
	defaultFields := sql_adaptor.FieldParseValidatorFromStruct(reflection)
	return sql_adaptor.NewSqlAdaptor(fieldMappings, defaultFields, matchers), nil
}

// keyValueMatcher is a custom matcher for tags[x].
func keyValueMatcher(ex *parser.Expression) (*sql_adaptor.SqlResponse, error) {
	reg, err := regexp.Compile(KeyValueRegex)
	if err != nil {
		return nil, err
	}
	slice := reg.FindStringSubmatch(ex.Field)
	if slice == nil {
		return nil, errors.New("didn't match regex expression")
	}
	if len(slice) < 3 {
		return nil, errors.New("regex match slice is too short")
	}
	if strings.ToLower(slice[1]) != "tags" {
		return nil, errors.New("expected tags as field name")
	}

	// We need to handle the % comparator differently since it isn't implicitly supported in SQL.
	defaultMatch := sql_adaptor.DefaultMatcher(&parser.Expression{
		Field:      "value",
		Comparator: ex.Comparator,
		Value:      ex.Value,
	})
	rawSnippet := defaultMatch.Raw
	if len(defaultMatch.Values) != 1 {
		return nil, errors.New("unexpected number of values from matcher")
	}
	value := defaultMatch.Values[0]
	sq := sql_adaptor.SqlResponse{
		Raw:    fmt.Sprintf("id IN (SELECT model_id FROM %s WHERE key=? AND %s AND deleted_at is NULL)", slice[1], rawSnippet),
		Values: []string{slice[2], value},
	}
	return &sq, nil
}
