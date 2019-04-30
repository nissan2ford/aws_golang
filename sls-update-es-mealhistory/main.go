package main

import (
	"fmt"

//	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"


//	"gopkg.in/olivere/elastic.v3"
//	"github.com/edoardo849/apex-aws-signer"
)

// define elasticsearch setting
const (
	// dynamoDB Table Name
	ddbTablename = "meal_history"

	esUrl = "https://search-es-mealhistory-yqcvq335bq5okinup2ohti24gy.ap-northeast-1.es.amazonaws.com"
//	esIndex = 
)

// define mealhistory type
type MDdata struct {
	MealTime	string	`json:"MealTime"`	// primary pertition key
	Date		int	`json:"Date"`		// primary sort key
	Day_of_week	string	`json:"Day_of_week"`
	MealMethod	string	`json:"MealMethod"`
}

type Response struct {
    Message string `json:"message"`
    Ok      bool   `json:"ok"`
}


//func PutES()(Response, error) {
//		transport := signer.NewTransport(session.New(&aws.Config{Region:aws.String(rec.AWSRegion)}), elasticsearchservice.ServiceName)

//		httpClient := &http.Client{
//			Transport: transport,
//		}

		// Elasticsearchに接続
//		client, err := elastic.NewClient(
//			elastic.SetSniff(false),
//			elastic.SetURL(esUrl), // AWS Elasticsearchのエンドポイント
//			elastic.SetScheme("https"),
//			elastic.SetHttpClient(httpClient),
//		)
//}

func Handler() (Response, error) {
	// 処理部分
	msg, err := Scanmealhistory()

	return msg,err
}


func main() {
    lambda.Start(Handler)
}