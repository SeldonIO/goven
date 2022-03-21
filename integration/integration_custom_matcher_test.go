package integration

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	. "github.com/onsi/gomega"
	"github.com/seldonio/goven/example"
	"log"
	"testing"
)


type testRigModel struct {
	pg      *embeddedpostgres.EmbeddedPostgres
	modelDAO *example.ModelDAO
}

func newTestRigModel() (*testRigModel, error) {
	db, pg, err := setupDB()
	if err != nil {
		return nil, err
	}
	// Create Model table
	err = db.AutoMigrate(example.Model{}, example.Tag{})
	if err != nil {
		return nil, err
	}
	dao, err := example.NewModelDAO(db)
	if err != nil {
		return nil, err
	}
	return &testRigModel{
		pg:      pg,
		modelDAO: dao,
	}, nil
}

func (t *testRigModel) cleanup() {
	err := t.pg.Stop()
	if err != nil {
		log.Print(err)
	}
}

func TestSqlAdaptorModel(t *testing.T) {
	g := NewGomegaWithT(t)
	rig, err := newTestRigModel()
	defer rig.cleanup()
	g.Expect(err).To(BeNil())
	// Setup entries
	model1Tags := []example.Tag{
		example.Tag{
			Key:     "auto_created",
			Value:   "true",
		},
	}
	err = rig.modelDAO.CreateModel(&example.Model{
		Name: "model1",
		Tags: model1Tags,
	})
	g.Expect(err).To(BeNil())
	model2Tags := []example.Tag{
		example.Tag{
			Key:     "auto_created",
			Value:   "false",
		},
	}
	err = rig.modelDAO.CreateModel(&example.Model{
		Name: "model2",
		Tags: model2Tags,
	})
	g.Expect(err).To(BeNil())
	t.Run("test simple successful query", func(t *testing.T) {
		result, err := rig.modelDAO.MakeQuery("name=model1")
		g.Expect(err).To(BeNil())
		g.Expect(len(result)).To(Equal(1))
		g.Expect(result[0].Name).To(Equal("model1"))
		g.Expect(len(result[0].Tags)).To(Equal(1))
		g.Expect(result[0].Tags[0].Key).To(Equal("auto_created"))
	})
	t.Run("test model tags true", func(t *testing.T) {
		result, err := rig.modelDAO.MakeQuery(`tags[auto_created]="true"`)
		g.Expect(err).To(BeNil())
		g.Expect(len(result)).To(Equal(1))
		g.Expect(result[0].Name).To(Equal("model1"))
		g.Expect(len(result[0].Tags)).To(Equal(1))
		g.Expect(result[0].Tags[0].Value).To(Equal("true"))
	})
	t.Run("test model tags false", func(t *testing.T) {
		result, err := rig.modelDAO.MakeQuery(`tags[auto_created]="false"`)
		g.Expect(err).To(BeNil())
		g.Expect(len(result)).To(Equal(1))
		g.Expect(result[0].Name).To(Equal("model2"))
		g.Expect(len(result[0].Tags)).To(Equal(1))
		g.Expect(result[0].Tags[0].Value).To(Equal("false"))
	})
}
