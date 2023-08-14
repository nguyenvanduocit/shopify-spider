package pkg

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func NewMongoDb() (*mongo.Client, func(), error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, nil, errors.New("MONGO_URI is empty")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to connect to mongo")
	}

	cleanup := func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}

	return client, cleanup, nil
}

func GetAppByUrl(collection *mongo.Collection, url string) (*App, error) {
	filter := bson.M{
		"url": url,
	}

	var app App

	if err := collection.FindOne(context.Background(), filter).Decode(&app); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get app by handler")
	}

	return &app, nil
}

func CreateApplication(collection *mongo.Collection, app *App) error {
	if _, err := collection.InsertOne(context.Background(), app); err != nil {
		return errors.Wrap(err, "failed to insert app")
	}

	return nil
}
