package main

import (
//	"fmt"

	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/dynamodb"
//	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"


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

//type Response struct {
//    Message string `json:"message"`
//    Ok      bool   `json:"ok"`
//}


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

	return string(b), nil
}

func Handler(input *MHdata) (string, error) {
	// 処理部分
	msg, err := PutES(input)

	return msg,err
}


func main() {
    lambda.Start(Handler)
}