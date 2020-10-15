package schema

import "github.com/facebook/ent"

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return nil
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return nil
}
