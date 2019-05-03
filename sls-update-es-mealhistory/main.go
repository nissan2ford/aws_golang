package main

import (
	"fmt"

	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"


	"log"
	"net/http"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"gopkg.in/olivere/elastic.v3"
	"github.com/edoardo849/apex-aws-signer"
)

// define elasticsearch setting
const (
	// dynamoDB Table Name
	ddbTablename = "meal_history"

	esUrl = "https://search-es-mealhistory-yqcvq335bq5okinup2ohti24gy.ap-northeast-1.es.amazonaws.com"
	esIndex = "mhindex"  // Index name
	esType = "mealhistory" // Index group name

	// AWS Region
	region = "ap-northeast-1"
)

// define mealhistory type
type MHdata struct {
	MealTime	string	`json:"MealTime"`	// primary pertition key
	Date		int	`json:"Date"`		// primary sort key
	Day_of_week	string	`json:"Day_of_week"`
	MealMethod	string	`json:"MealMethod"`
}

func PutES(mhdata *MHdata)(string, error) {
		transport := signer.NewTransport(session.New(&aws.Config{Region:aws.String(region)}), elasticsearchservice.ServiceName)

		httpClient := &http.Client{
			Transport: transport,
		}

		// Elasticsearchに接続
		client, err := elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL(esUrl), // AWS Elasticsearchのエンドポイント
			elastic.SetScheme("https"),
			elastic.SetHttpClient(httpClient),
		)

		if err != nil {
			panic(err)
		}

		// Create an index.
		indexName := esIndex

		// Index a recode
		msg, err := putDataIntoES(client,indexName,mhdata)

	return msg, nil
}

// Index a recode of MHdata
func putDataIntoES(c *elastic.Client, indexName string,mhdata *MHdata)(string, error) {

	log.Printf("mhdata: %s\n", mhdata)

	putdata, err := c.Index().
		Index(indexName).
		Type(esType).
		BodyJson(mhdata).
		Do()

	if err != nil {
		panic(err)
	}

	log.Printf("Indexed mhdata %s to index %s, type %s\n",putdata.Id,putdata.Index, putdata.Type)

	// Convert JSON
	b, err := json.Marshal(mhdata)

	if err != nil {
		panic(err)
	}

	return string(b), nil
}

func UnmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {

	dbAttrMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attribute {

		var dbAttr dynamodb.AttributeValue

		bytes, marshalErr := v.MarshalJSON(); if marshalErr != nil {
			return marshalErr
		}

		json.Unmarshal(bytes, &dbAttr)
		dbAttrMap[k] = &dbAttr
	}

	return dynamodbattribute.UnmarshalMap(dbAttrMap, out)
}

func handleRequest(ctx context.Context, e events.DynamoDBEvent) (string, error) {

	// initialize var
	mhdata := MHdata{}

	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		// Print new values for attributes of type String
		for name, value := range record.Change.NewImage {
			if value.DataType() == events.DataTypeString {
				fmt.Printf("Attribute name: %s, value: %s\n", name, value.String())

			}
			
			err := UnmarshalStreamImage(record.Change.NewImage, &mhdata)
			if err != nil {
				panic(err)
			}
		}
	}

	// Put elasticsearch
	msg, err := PutES(&mhdata)

	if err != nil {
		panic(err)
	}

	return msg,err

}

func main() {
    lambda.Start(handleRequest)
}