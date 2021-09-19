package repository

import (
	"context"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Upsert stores user in the database.
func (r *UsersMongo) Upsert(u model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	var (
		filter = bson.M{"id": u.ID}
		update = bson.M{ "$set": u }
	)

	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	return err
}

// UpdateByID partially updates user in the database by given `id`.
func (r *UsersMongo) UpdateByID(id string, set map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	var (
		filter = bson.M{"id": id}
		update = bson.M{"$set": set}
	)

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

// GetByID retrieves user from the collection by given `userID`.
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

// GetByQuery retrieves user from the collection by given `query`.
func (r *UsersMongo) GetByQuery(query map[string]interface{}) (*model.User, error) {
	var (
		user *model.User
	)

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	err := r.collection.FindOne(ctx, query).Decode(&user)

	return user, err
}

// ListByQuery retrieves users from the collection by given `query`.
func (r *UsersMongo) ListByQuery(query map[string]interface{}) ([]*model.User, error) {
	var (
		users []*model.User
	)

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	cursor.All(ctx, &users)

	return users, err
}
