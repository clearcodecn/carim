package schema

import "github.com/facebook/ent"

// GroupUser holds the schema definition for the GroupUser entity.
type GroupUser struct {
	ent.Schema
}

// Fields of the GroupUser.
func (GroupUser) Fields() []ent.Field {
	return nil
}

// Edges of the GroupUser.
func (GroupUser) Edges() []ent.Edge {
	return nil
}
