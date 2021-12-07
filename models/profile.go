package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Profile struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Bio       string             `json:"bio"`
	BirthDate time.Time          `json:"birth_date" bson:"-"`
	YOB       int                `bson:"year_of_birth"`
	Name      string             `json:"name"`
	Photos    []struct {
		URL       string `json:"url"`
		MediaType string `json:"media_type" bson:"media_type"`
		Extension string `json:"extension"`
	} `json:"photos"`
	LastModified   time.Time `json:"last_modified" bson:"last_modified"`
	RecentlyActive bool      `json:"recently_active"`
}
