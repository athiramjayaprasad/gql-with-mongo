package job

import (
	"context"
	"fmt"
	"time"

	"github.com/athiramjayaprasad/gql-with-mongo/database"
	"github.com/athiramjayaprasad/gql-with-mongo/graph/model"
	"go.mongodb.org/mongo-driver/bson"
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
	query := database.NewQuery("Job", database.WithSelectFields([]string{"_id", "title", "description", "company", "url"}))
	cursor, err :=j.mongoJobCollection.List(query, ctx)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
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
		ID: jobId.(string),
		Title: input.Title,                                        
		Description: input.Description,
		URL: input.URL,
		Company: input.Company,
	}

	return &job, nil
}
                                             