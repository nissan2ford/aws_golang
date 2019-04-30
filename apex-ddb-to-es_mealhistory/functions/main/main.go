package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"


//	"gopkg.in/olivere/elastic.v3"
//	"github.com/edoardo849/apex-aws-signer"
)

// define elasticsearch setting
const (
	// dynamoDB Table Name
	ddbTablename = "mealhistory"

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

func Scanmealhistory()(Response, error) {
	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)

	// Query
	queryParams := &dynamodb.QueryInput {
		TableName: aws.String(ddbTablename),
		KeyConditionExpression: aws.String("#MealTime=:mealtime"),
		ExpressionAttributeNames: map[string]*string {
			"#MealTime": aws.String("mealtime"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
			":mealtime": {
				S: aws.String("MOR"),
//				S: aws.String("LUN"),
//				S: aws.String("DIN"),
			},
		},
	}

	queryItem, queryErr := svc.Query(queryParams)
	if queryErr != nil {
		panic(queryErr)
	}

	fmt.Println(queryItem)

	return Response{
		Message: fmt.Sprintln(queryItem),
		Ok:      true,
	}, nil
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