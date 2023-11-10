package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/athiramjayaprasad/gql-with-mongo/database"
	"github.com/athiramjayaprasad/gql-with-mongo/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewJobDao() *JobDao{
   return &JobDao{
	mongoJobCollection: database.NewMongoCollection(database.MongoClient, database.DB, 30*time.Second),
   }  
}

type JobDao struct {
	mongoJobCollection database.MongoCollection
}

func (j *JobDao) List() ([]*model.JobListing, error) {
	var results []*model.JobListing
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	query := database.NewQuery("Job", database.WithSelectFields([]string{"title", "description", "company", "url"}))
	cursor, err :=j.mongoJobCollection.List(query, ctx)
	if err != nil {
		return nil, err
	}
	var resultBson []bson.M

	err = cursor.All(context.TODO(), &resultBson)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(resultBson)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(jsonData, &results)
	fmt.Println(results)
	return results, err
}

func (j *JobDao) Save(input model.CreateJobListingInput) (*model.JobListing, error) {
	insertData := bson.M{
		"title" : input.Title,
		"description": input.Description,
		"url" : input.URL,
		"company": input.Company,
	}
	jobId, err := j.mongoJobCollection.Save(insertData, "Job")
	if err != nil {
		return nil, err
	}
	job := model.JobListing{
		ID: jobId,
		Title: input.Title,                                        
		Description: input.Description,
		URL: input.URL,
		Company: input.Company,
	}

	return &job, nil
}

func (j *JobDao) Retrieve(id string) (*model.JobListing, error) {
	var result *model.JobListing
	var resultBson bson.M
	_id, _ := primitive.ObjectIDFromHex(id)

	query := database.NewQuery("Job",database.WithFilter([]database.Filter{{Key: "_id", Value: _id, Condition: database.Equals}}), database.WithSelectFields([]string{"_id", "title", "description", "company", "url"}))
    err :=j.mongoJobCollection.Retrieve(query, "Job").Decode(&resultBson)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the Id %s\n", id)
	}    
	if err != nil {
		panic(err)                   
	}    
	jsonData, err := json.Marshal(resultBson)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(jsonData, &result)
	return result, err 

}