package example

import (
	"context"
	"errors"
	"fmt"
	"github.com/seldonio/goven/parser"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/seldonio/goven/sql_adaptor"
	"gorm.io/gorm"
)

const (
	KeyValueRegex = `(.+)\[(.+)\]`
)

type Model struct {
	gorm.Model
	Name         string
	Version      string
	CreatedAt    time.Time
	Tags         []Tag      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Tag struct {
	gorm.Model
	Key             string
	Value           string
	ModelID 		uint
}

type ModelDAO struct {
	db           *gorm.DB
	queryAdaptor *sql_adaptor.SqlAdaptor
}

func NewModelDAO(db *gorm.DB) (*ModelDAO, error) {
	adaptor, err := createModelAdaptor()
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
	query = query.Model(Model{}).Where(queryResp.Raw, sql_adaptor.StringSliceToInterfaceSlice(queryResp.Values)...)
	err = query.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func createModelAdaptor() (*sql_adaptor.SqlAdaptor, error) {
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
	sq := sql_adaptor.SqlResponse{
		Raw:    fmt.Sprintf("id IN (SELECT model_id FROM %s WHERE key=? AND value%s?)", slice[1], ex.Comparator),
		Values: []string{slice[2], ex.Value},
	}
	return &sq, nil
}
