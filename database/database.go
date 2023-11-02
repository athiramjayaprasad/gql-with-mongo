package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionString string = "mongodb+srv://athiramjayaprasad23:Athira@cluster0.2osixpd.mongodb.net/?retryWrites=true&w=majority"
var MongoClient *mongo.Client

func init() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	MongoClient, _ = mongo.Connect(context.TODO(), opts)
}

type MongoModel interface {
	GetID() string
	CollectionName() string
}

type MongoCollection[M MongoModel] struct {
	mongoClient *mongo.Client
	db          string
	QueryTimeout time.Duration
}


func NewMongoCollection[M MongoModel](client *mongo.Client, db string) MongoCollection[M] {
	return MongoCollection[M]{
		mongoClient: client,
		db:          db,
	}
}

func (m MongoCollection[M]) List(query Query) ([]M, error) {
	var model M
	var results []M
	ctx, cancel := context.WithTimeout(context.Background(), m.QueryTimeout)
	defer cancel()
	filters, opts := query.MongoQuery()
	dbCollection := m.mongoClient.Database(m.db).Collection(model.CollectionName())
	data, err := dbCollection.Find(ctx, filters, opts)
	if err != nil {
		return results, err 
	}
	err = data.All(ctx, results)
	return results, err
}
