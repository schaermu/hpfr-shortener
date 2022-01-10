package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortURL struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at,omitempty"`
	ShortCode     string             `bson:"short_code,omitempty"`
	TargetURL     string             `bson:"target_url,omitempty"`
	RedirectCount int                `bson:"redirect_count,omitempty"`
}
