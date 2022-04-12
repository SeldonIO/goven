// Package integration contains integration tests using the example DAOs.
package integration

import (
	"log"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	. "github.com/onsi/gomega"
	"github.com/seldonio/goven/example"
)

var (
	model1 = &example.Model{
		Name: "model1",
		Tags: []example.Tag{
			{
				Key:   "auto_created",
				Value: "true",
			},
		},
	}
	model2 = &example.Model{
		Name: "model2",
		Tags: []example.Tag{
			{
				Key:   "auto_created",
				Value: "false",
			},
		},
	}
	deployment1 = &example.Model{
		Name: "deployment1",
		Tags: []example.Tag{
			{
				Key:   "tag",
				Value: "test_partial1",
			},
		},
	}
	deployment2 = &example.Model{
		Name: "deployment2",
		Tags: []example.Tag{
			{
				Key:   "tag",
				Value: "test_partial2",
			},
		},
	}
)

type testRigModel struct {
	pg       *embeddedpostgres.EmbeddedPostgres
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
		pg:       pg,
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
	for _, model := range []*example.Model{model1, model2, deployment1, deployment2} {
		err = rig.modelDAO.CreateModel(model)
		g.Expect(err).To(BeNil())
	}
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
	t.Run("test partial string match", func(t *testing.T) {
		result, err := rig.modelDAO.MakeQuery(`name%"model"`)
		g.Expect(err).To(BeNil())
		g.Expect(len(result)).To(Equal(2))
	})
	t.Run("test partial string tags", func(t *testing.T) {
		result, err := rig.modelDAO.MakeQuery(`tags[tag] % partial`)
		g.Expect(err).To(BeNil())
		g.Expect(len(result)).To(Equal(2))
	})
}
