package schema

import "github.com/facebook/ent"

// ChatList holds the schema definition for the ChatList entity.
type ChatList struct {
	ent.Schema
}

// Fields of the ChatList.
func (ChatList) Fields() []ent.Field {
	return nil
}

// Edges of the ChatList.
func (ChatList) Edges() []ent.Edge {
	return nil
}
