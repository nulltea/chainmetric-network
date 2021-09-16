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

// GetAll retrieves all model.Subscription from the collection.
func (r *SubscriptionsMongo) GetAll() ([]model.SubscriptionTicket, error) {
	var (
		subs []model.SubscriptionTicket
		ctx, cancel = context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	)


	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	cursor.All(ctx, &subs)

	return subs, err
}

// Insert stores model.Subscription in the database.
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
