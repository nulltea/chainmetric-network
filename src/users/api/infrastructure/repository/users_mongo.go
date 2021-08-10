package repository

import (
	"context"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-contracts/model/user"
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

// Store stores user in the database.
func (r *UsersMongo) Store(u user.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	_, err := r.collection.InsertOne(ctx, u)

	return err
}

// GetByID retrieves user from the collection by given `userID`.
func (r *UsersMongo) GetByID(userID string) (*user.User, error) {
	var (
		user *user.User
		filter = bson.M{"id": userID}
	)

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	err := r.collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}
