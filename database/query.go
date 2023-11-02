package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Filter struct {
	Key       string
	Condition Operator
	Value     interface{}
}

type Update struct {
	Key   string
	Value interface{}
}

type Direction int

const (
	Asc  Direction = 1
	Desc Direction = -1
)

type Sort struct {
	Key       string
	Direction Direction
}

type Operator string

const (
	Equals              Operator = "$eq"
	In                  Operator = "$in"
	NotIn               Operator = "$nin"
	GreaterThan         Operator = "$gt"
	GreaterThanOrEquals Operator = "$gte"
	LessThan            Operator = "$lt"
	LessThanOrEquals    Operator = "$lte"
	NotEquals           Operator = "$ne"
	Or                  Operator = "$or"
	And                 Operator = "$and"
	Not                 Operator = "$not"
	Regex               Operator = "$regex"
	Exists              Operator = "$exists"
	Set                 Operator = "$set"
)

type Query struct {
	offset     int
	limit      int
	collection string
	fields     []string
	filterOpts []Filter
	sortOpts   []Sort
}

type QueryOption func(q *Query)

func WithFilter(filters []Filter) func(q *Query) {
	return func(q *Query) {
		q.filterOpts = filters
	}
}

func WithOffset(offset int) func(q *Query) {
	return func(q *Query) {
		q.offset = offset
	}
}

func WithLimit(limit int) func(q *Query) {
	return func(q *Query) {
		q.limit = limit
	}
}

func WithSelectFields(fields []string) func(q *Query) {
	return func(q *Query) {
		q.fields = fields
	}
}

func WithSort(sorts []Sort) func(q *Query) {
	return func(q *Query) {
		q.sortOpts = sorts
	}
}

func NewQuery(opts ...QueryOption) Query {
	q := &Query{}
	for _, opt := range opts {
		opt(q)
	}
	return *q
}

func (q Query) MongoQuery() (primitive.D, *options.FindOptions) {
	filters := bson.D{}
	filter := bson.E{}

	for _, filterOpt := range q.filterOpts {
		// assume equals in case filter condition not found
		if filterOpt.Condition == "" {
      
		} else {
			filter = bson.E{Key: filterOpt.Key, Value: bson.M{string(filterOpt.Condition): filterOpt.Value}}
		}
		filters = append(filters, filter)
	}

	sortOrder := map[string]int{}   
	for _, sortOpt := range q.sortOpts {
		sortOrder[sortOpt.Key] = int(sortOpt.Direction)
	}

	projections := map[string]int{}
	for _, field := range q.fields {
		projections[field] = 1
	}
    
	opts := options.Find()
	opts.SetLimit(int64(q.limit))
	opts.SetSkip(int64(q.offset))
	opts.Sort = sortOrder
	opts.SetProjection(projections)

	return filters, opts 
}
