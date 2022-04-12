package example

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/seldonio/goven/sql_adaptor"
	"gorm.io/gorm"
)

// User represents an simple example database schema.
type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	CreatedAt    time.Time
}

// UserDAO is an example DAO for user data.
type UserDAO struct {
	db           *gorm.DB
	queryAdaptor *sql_adaptor.SqlAdaptor
}

// NewUserDAO returns a UserDAO.
func NewUserDAO(db *gorm.DB) (*UserDAO, error) {
	reflection := reflect.ValueOf(&User{})
	adaptor := sql_adaptor.NewDefaultAdaptorFromStruct(reflection)
	return &UserDAO{
		db:           db,
		queryAdaptor: adaptor,
	}, nil
}

// CreateUser commits the provided user to the database.
func (u *UserDAO) CreateUser(user *User) error {
	ctx := context.Background()
	tx := u.db.Begin().WithContext(ctx)
	err := tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// MakeQuery takes a goven query and performs it against the user database.
func (u *UserDAO) MakeQuery(q string) ([]User, error) {
	var users []User
	ctx := context.Background()
	query := u.db.WithContext(ctx)
	queryResp, err := u.queryAdaptor.Parse(q)
	if err != nil {
		return nil, err
	}
	query = query.Model(User{}).Where(queryResp.Raw, sql_adaptor.StringSliceToInterfaceSlice(queryResp.Values)...)
	err = query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
