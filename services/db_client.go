package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DbCallTimeoutSec is maximum time to wait for db response after which call is canceled
const DbCallTimeoutSec = 30

// InsertManyBatchMaxSize is a maximum number of elements to be written in a single insert many batch operation
const InsertManyBatchMaxSize = 100_000

// TokenLength is number of bytes in token string
const TokenLength = 128

// ThisApiName is name of this microservices,
// all data refering to this microservice are written under this name in database
const ThisApiName = "go_basilisk"

// Collection represents mongo db collection name
type Collection string

// Collections is pseudo enum to access mongo db collection
var Collections = struct {
	Tokens string
}{
	Tokens: "tokens",
}

// TokenDB represents token document in database
type TokenDB struct {
	APIName string `bson:"api_name"`
	Token   string `bson:"token"`
}

// MongoClient handles connection to database
type MongoClient struct {
	inner mongo.Database
}

// NewMongoClient creates instance of MongoClient
// Connects to mongo db database and returns pointer to new instance structure of MongoClient
func NewMongoClient(uri, database string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &MongoClient{inner: *client.Database(database)}, nil
}

// GetToken returns token
func (mc *MongoClient) GetToken() (string, error) {
	c := mc.inner.Collection(Collections.Tokens)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*DbCallTimeoutSec)
	defer cancel()
	query := bson.M{"api_name": ThisApiName}
	var t TokenDB
	err := c.FindOne(ctx, query).Decode(&t)
	if err != nil {
		log.Printf("error in getToken method, err: %s\n", err)
		return "", err
	}
	return t.Token, nil
}

// SetToken sets fresh token in database
func (mc *MongoClient) SetToken() error {
	c := mc.inner.Collection(Collections.Tokens)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*DbCallTimeoutSec)
	defer cancel()
	ts, err := RandToken(TokenLength)
	if err != nil {
		return err
	}
	var t TokenDB
	filter := bson.M{"api_name": ThisApiName, "token": ts}
	err = c.FindOne(ctx, filter).Decode(&t)
	if err != nil {
		log.Printf("error when looking for token: %s\n", err)
	}
	token := bson.M{"api_name": ThisApiName, "token": ts}
	if t.APIName == "" && t.Token == "" {
		res, err := c.InsertOne(ctx, token)
		if err != nil {
			return err
		}
		log.Printf("upserted id %v\n", res.InsertedID)
		return nil
	}
	update := bson.M{"$set": token}
	res, err := c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	log.Printf("upserted %v\n, modified %v\n", res.ModifiedCount, res.ModifiedCount)
	return nil

}
