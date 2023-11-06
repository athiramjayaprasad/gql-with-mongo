package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionString string = "mongodb+srv://athiramjayaprasad23:Athira@cluster0.2osixpd.mongodb.net/?retryWrites=true&w=majority"
var MongoClient *mongo.Client
var DB string = "gql_example"

func init() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)
	var err error
	// Create a new client and connect to the server
	MongoClient, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		
	}
}

type MongoCollection struct {
	mongoClient *mongo.Client
	db          string
	QueryTimeout time.Duration
}


func NewMongoCollection(client *mongo.Client, db string, queryTimeout time.Duration) MongoCollection{
	return MongoCollection{
		mongoClient: client,
		db:          db,
		QueryTimeout: queryTimeout,
	}
}

func (m MongoCollection) List(query Query, ctx context.Context) (*mongo.Cursor, error) {
	filters, opts := query.MongoQuery()
	dbCollection := m.mongoClient.Database(m.db).Collection(query.collection)
	data, err := dbCollection.Find(ctx, filters, opts)
	if err != nil {
		return data, err 
	}
	
	return data, err
}

func (m MongoCollection) Save (obj interface{}, collectionName string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeout)
	defer cancel()
	dbCollection := m.mongoClient.Database(m.db).Collection(collectionName)
	result, err := dbCollection.InsertOne(ctx, obj)
	if err != nil {             
		return obj, err
	}
    
	id := result.InsertedID.(primitive.ObjectID).String()
	return id, nil
    
}