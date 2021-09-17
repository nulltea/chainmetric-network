package repository

import (
	"context"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SubscriptionsMongo defines subscriptions repository for MongoDB database.
type SubscriptionsMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewSubscriptionsMongo constructs new SubscriptionsMongo repository instance.
func NewSubscriptionsMongo(client *mongo.Client) *SubscriptionsMongo {
	return &SubscriptionsMongo{
		client:     client,
		collection: client.Database(viper.GetString("mongo_database")).Collection("subscriptions"),
	}
}

// GetBySubID retrieves subscription tickets from the collection by given subscription id.
func (r *SubscriptionsMongo) GetBySubID(subID string) ([]model.SubscriptionTicket, error) {
	var (
		results     []model.SubscriptionTicket
		ctx, cancel = context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	)


	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{
		"sub_id": subID,
	})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, err
}

// Insert stores subscription ticket in the database.
func (r *SubscriptionsMongo) Insert(tickets ...model.SubscriptionTicket) error {
	var (
		docs []interface{}
		ctx, cancel = context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	)

	defer cancel()

	for i := range tickets {
		docs = append(docs, tickets[i])
	}

	_, err := r.collection.InsertMany(ctx, docs)

	return err
}

// DeleteByIDs removes subscription ticket from the database by given `ids`.
func (r *SubscriptionsMongo) DeleteByIDs(ids ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))

	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{
		"id": bson.M{
			"$in": ids,
		},
	})

	return err
}
