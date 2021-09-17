package repository

import (
	"context"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EventConcernsMongo defines subscriptions repository for MongoDB database.
type EventConcernsMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewEventConcernsMongo constructs new EventConcernsMongo repository instance.
func NewEventConcernsMongo(client *mongo.Client) *EventConcernsMongo {
	return &EventConcernsMongo{
		client:     client,
		collection: client.Database(viper.GetString("mongo_database")).Collection("event_concerns"),
	}
}

// GetAll retrieves all model.EventConcern from the collection.
func (r *EventConcernsMongo) GetAll() ([]model.EventConcern, error) {
	var (
		results     []model.EventConcern
		ctx, cancel = context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	)


	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var (
			record model.EventConcern
			topic = model.EventTopic(cursor.Current.Lookup("event_topic").String())
		)

		switch topic {
		case model.RequirementsViolationTopic:
			record = new(model.RequirementsViolationConcern)
		default:
			core.Logrus.WithField("topic", topic).
				Warn("unknown topic: document cannot be decoded")
			break
		}

		if err = cursor.Decode(record); err != nil {
			results = append(results, record)
		}
	}

	return results, err
}

// Insert stores model.EventConcern in the database.
func (r *EventConcernsMongo) Insert(concerns ...model.EventConcern) error {
	var (
		docs []interface{}
		ctx, cancel = context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	)

	defer cancel()

	for i := range concerns {
		docs = append(docs, concerns[i])
	}

	_, err := r.collection.InsertMany(ctx, docs)

	return err
}

// DeleteByIDs removes model.EventConcern from the database by given `ids`.
func (r *EventConcernsMongo) DeleteByIDs(ids ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))

	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{
		"id": bson.M{
			"$in": ids,
		},
	})

	return err
}
