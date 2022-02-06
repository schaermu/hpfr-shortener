package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortURL struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	ShortCode string             `bson:"short_code,omitempty"`
	TargetURL string             `bson:"target_url,omitempty"`
	Hits      []ShortURLHit      `bson:"hits,omitempty"`
}

type ShortURLHit struct {
	CreatedAt    time.Time `bson:"created_at,omitempty"`
	UAFamily     string    `bson:"ua_family,omitempty"`
	UAMajor      string    `bson:"ua_major,omitempty"`
	UAMinor      string    `bson:"ua_minor,omitempty"`
	OS           string    `bson:"os_name,omitempty"`
	OSMajor      string    `bson:"os_major,omitempty"`
	OSMinor      string    `bson:"os_minor,omitempty"`
	DeviceFamily string    `bson:"device_family,omitempty"`
	Country      string    `bson:"country,omitempty"`
	CountryCode  string    `bson:"country_code,omitempty"`
}
