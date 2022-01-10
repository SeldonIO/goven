package integration

import (
	"log"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/seldonio/goven/example"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type testRig struct {
	pg      *embeddedpostgres.EmbeddedPostgres
	userDAO *example.UserDAO
}

func newTestRig() (*testRig, error) {
	// Create Postgres db
	config := embeddedpostgres.DefaultConfig().Port(9876)
	pg := embeddedpostgres.NewDatabase(config)
	err := pg.Start()
	if err != nil {
		return nil, err
	}

	// Connect gorm to db
	dsn := "host=localhost port=9876 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create User table
	err = db.AutoMigrate(example.User{})
	if err != nil {
		return nil, err
	}

	dao, err := example.NewUserDAO(db)
	if err != nil {
		return nil, err
	}
	return &testRig{
		pg:      pg,
		userDAO: dao,
	}, nil
}

func (t *testRig) cleanup() {
	err := t.pg.Stop()
	if err != nil {
		log.Print(err)
	}
}

func TestSqlAdaptor(t *testing.T) {
	g := NewGomegaWithT(t)
	rig, err := newTestRig()
	defer rig.cleanup()
	g.Expect(err).To(BeNil())

	// Setup entries
	err = rig.userDAO.CreateUser(&example.User{
		Name: "",
		Age:  10,
	})
	g.Expect(err).To(BeNil())
	err = rig.userDAO.CreateUser(&example.User{
		Name: "dom",
		Age:  12,
	})
	g.Expect(err).To(BeNil())
	err = rig.userDAO.CreateUser(&example.User{
		Name: "dom",
		Age:  9,
	})
	g.Expect(err).To(BeNil())
	t.Run("test simple successful query", func(t *testing.T) {
		result, err := rig.userDAO.MakeQuery("name=dom AND age>11")
		g.Expect(err).To(BeNil())
		g.Expect(len(result)).To(Equal(1))
		g.Expect(result[0].Name).To(Equal("dom"))
		g.Expect(result[0].Age).To(Equal(uint8(12)))
	})
	t.Run("test empty string query", func(t *testing.T) {
		result, err := rig.userDAO.MakeQuery("name=\"\"")
		g.Expect(err).To(BeNil())
		g.Expect(len(result)).To(Equal(1))
		g.Expect(result[0].Name).To(Equal(""))
		g.Expect(result[0].Age).To(Equal(uint8(10)))
	})
}
