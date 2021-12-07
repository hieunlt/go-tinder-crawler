package infrastructure

import (
	"context"
	"go-tinder-crawler/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
	"time"
)

type DBConn struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewDBConn(uri string) (*DBConn, error) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database(os.Getenv("DB_NAME")).Collection("profiles")
	log.Println("DB connection established")
	return &DBConn{client: client, collection: collection}, nil
}

func (c *DBConn) UpsertProfiles(profiles []models.Profile) error {
	var upsertResult UpsertResult
	for _, profile := range profiles {
		if len(strings.TrimSpace(profile.Bio)) == 0 {
			log.Printf("%s has empty bio, skipping\n", profile.ID)
			continue
		}
		filter := bson.D{{"_id", profile.ID}}
		replaceOptions := options.Replace().SetUpsert(true)
		result, err := c.collection.ReplaceOne(context.TODO(), filter, profile, replaceOptions)
		if err != nil {
			return err
		}
		log.Printf("%s upserted with result: %+v\n", profile.ID, *result)
		upsertResult.update(result)

	}
	log.Printf("%d profile(s) fetched with result: %+v\n", len(profiles), upsertResult)
	return nil
}

func (c *DBConn) Close() error {
	if err := c.client.Disconnect(context.TODO()); err != nil {
		return err
	}
	log.Println("DB connection closed")
	return nil
}

type UpsertResult struct {
	MatchedCount  int64 // The number of documents matched by the filter.
	ModifiedCount int64 // The number of documents modified by the operation.
	UpsertedCount int64 // The number of documents upserted by the operation.
}

func (r *UpsertResult) update(updateResult *mongo.UpdateResult) {
	r.UpsertedCount += updateResult.UpsertedCount
	r.ModifiedCount += updateResult.ModifiedCount
	r.MatchedCount += updateResult.MatchedCount
}
