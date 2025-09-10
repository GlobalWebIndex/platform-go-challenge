package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/google/uuid"
)

// Favourite represents a "user favourites an asset" relation.
type Favourite struct {
	ent.Schema
}

func (Favourite) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),

		field.UUID("user_id", uuid.UUID{}),
		field.UUID("asset_id", uuid.UUID{}),

		field.Time("created_at").
			Default(func() time.Time {
				return time.Now().UTC()
			}).
			Immutable(),
	}
}

func (Favourite) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("favourites").
			Field("user_id").
			Required().
			Unique(),

		edge.From("asset", Asset.Type).
			Ref("favourites").
			Field("asset_id").
			Required().
			Unique(),
	}
}

func (Favourite) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "asset_id").Unique(),
		index.Fields("user_id", "created_at", "id"),
	}
}
