package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	var fields = []ent.Field{
		field.String("nickname").Default(""),
		field.String("car_no").Unique(),
		field.String("avatar").Default(""),
		field.Enum("sex").Values("1", "2", "3"),
		field.String("city").Optional(),
		field.String("province").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.Time("deleted_at").Optional(),
	}
	return fields
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
