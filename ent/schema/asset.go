package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/google/uuid"
)

// Asset holds the schema definition for the Asset entity.
type Asset struct {
	ent.Schema
}

// Fields returns the Asset fields.
func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),

		// Controlled vocabulary. Ent generates a Go const set + validation.
		field.Enum("asset_type").Immutable().
			Values("chart", "insight", "audience"),

		// Required non-empty text.
		field.String("description").
			NotEmpty(),

		// JSON payload we store as Postgres JSONB.
		//
		// - Using map[string]any makes it easy to work with dynamic payloads.
		// - field.JSON tells Ent to JSON-encode/decode automatically in Go.
		// - SchemaType overrides the column type specifically for Postgres to "jsonb".
		field.JSON("payload", map[string]any{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}),

		field.Time("created_at").
			Default(func() time.Time {
				return time.Now().UTC()
			}).
			Immutable(),
	}
}

func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("favourites", Favourite.Type),
	}
}
