package mongodb

import (
	"RoomTgBot/internal/consts"
	"RoomTgBot/internal/types"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---------------- DB interactions -----------------------------

// Add one object to a specified collection
func AddOne[mongoObject types.MongoObject](ctx context.Context, client *mongo.Client, collectionName string, object *mongoObject) (*mongo.InsertOneResult, error) {
	collection := client.Database(consts.MongoDBName).Collection(collectionName)
	insertResult, err := collection.InsertOne(ctx, object)

	if err != nil {
		return insertResult, fmt.Errorf("unable to add new %s to MongoDB: %v", collectionName, err)
	}

	return insertResult, nil
}

// Get all objects from specified collection
func GetAll[mongoObject types.MongoObject](ctx context.Context, client *mongo.Client, collectionName string) ([]mongoObject, error) {
	collection := client.Database(consts.MongoDBName).Collection(collectionName)
	getError := func(err error) error {
		return fmt.Errorf("unable to get %s from MongoDB: %v", collectionName, err)
	}

	cursor, err := collection.Find(ctx, bson.D{})
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

func UpdateOne[mongoObject types.MongoObject](ctx context.Context, client *mongo.Client, collectionName string, object mongoObject) error {
	collection := client.Database(consts.MongoDBName).Collection(collectionName)
	filter := bson.M{"_id": object.MongoId()}
	update := bson.M{"$set": object}

    _, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("unable to update object due to : %v", err)
	}

	return nil
}

func UpdateAll[mongoObject types.MongoObject](ctx context.Context, client *mongo.Client, collectionName string, objects []mongoObject) error {
	collection := client.Database(consts.MongoDBName).Collection(collectionName)

    for _, elem := range objects {
		filter := bson.M{"_id": elem.MongoId()}
		update := bson.M{"$set": elem}

        _, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("unable to update object due to : %v", err)
		}
	}

	return nil
}

func Ping(client *mongo.Client) error {
	if client == nil {
		return fmt.Errorf("MongoDB client is nil")
	}

	return client.Ping(context.TODO(), nil)
}

func uri() (string, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		err := fmt.Errorf("'MONGODB_URI' is not set as environmental variable")
		return uri, err
	}

    return uri, nil
}

func NewClient() (*mongo.Client, error) {
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

func Disconnect(client *mongo.Client) {
	client.Disconnect(context.TODO())
}
