package repository

import (
	"context"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UsersMongo defines users repository for MongoDB database.
type UsersMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewUserMongo constructs new UsersMongo repository instance.
func NewUserMongo(client *mongo.Client) *UsersMongo {
	return &UsersMongo{
		client:     client,
		collection: client.Database(viper.GetString("mongo_database")).Collection("users"),
	}
}

// GetByID retrieves model from the collection by given `userID`.
func (r *UsersMongo) GetByID(userID string) (*model.User, error) {
	var (
		user *model.User
		filter = bson.M{"id": userID}
	)

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	err := r.collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}
