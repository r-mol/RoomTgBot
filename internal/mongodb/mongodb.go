package mongodb

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/types"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// ---------------- DB interactions -----------------------------

func AddOne[mongoObject types.MongoObject](ctx context.Context, client *mongo.Client, collectionName string, object mongoObject) (*mongo.InsertOneResult, error) {
	collection := client.Database(consts.MongoDBName).Collection(collectionName)
	insertResult, err := collection.InsertOne(ctx, object)

	if err != nil {
		return insertResult, fmt.Errorf("unable to add new %s to MongoDB: %v", collectionName, err)
	}

	return insertResult, nil
}

func GetAll[mongoObject types.MongoObject](ctx context.Context, client *mongo.Client, collectionName string) ([]mongoObject, error) {
	collection := client.Database(consts.MongoDBName).Collection("users")
	getError := func(err error) error {
		return fmt.Errorf("unable to get %s from MongoDB: %v", collectionName, err)
	}

	cursor, err := collection.Find(ctx, nil)
	if err != nil {
		return []mongoObject{}, getError(err)
	}

	users := []mongoObject{}

	for cursor.Next(context.TODO()) {
		var result mongoObject
		if err := cursor.Decode(&result); err != nil {
			return []mongoObject{}, getError(err)
		}

		users = append(users, result)
	}

	if err := cursor.Err(); err != nil {
		return []mongoObject{}, getError(err)
	}

	return users, nil
}

// ---------------- DB initialization -----------------------------

func init() {
	var err error
	mongoClient, err = newClient()

	if err != nil {
		panic(err)
	}

	err = Ping(mongoClient)
	if err != nil {
		panic(fmt.Errorf("Ping to MongoDB is unsuccessful: %v", err))
	}
}

func Ping(client *mongo.Client) error {
	if client == nil {
		return fmt.Errorf("MongoDB client is nil")
	}
	return client.Ping(context.TODO(), nil)
}

func Client() *mongo.Client {
	return mongoClient
}

func uri() (string, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		err := fmt.Errorf("'MONGODB_URI' is not set as environmental variable")
		return uri, err
	}
	return uri, nil
}

func newClient() (*mongo.Client, error) {
	uri, err := uri()
	if err != nil {
		return nil, err
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}
