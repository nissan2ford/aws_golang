package main

import (
	"fmt"

	"log"

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

//	esUrl = "https://search-es-mealhistory-yqcvq335bq5okinup2ohti24gy.ap-northeast-1.es.amazonaws.com"
//	esIndex = 
)

// define mealhistory type
type MHdata struct {
	MealTime	string	`json:"MealTime"`	// primary pertition key
	Date		int	`json:"Date"`		// primary sort key
	Day_of_week	string	`json:"Day_of_week"`
	MealMethod	string	`json:"MealMethod"`
}

func Scanmealhistory(mhdatas *[]MHdata)(string, error) {
	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess,aws.NewConfig().WithRegion("ap-northeast-1"))

	// Query
	queryParams := &dynamodb.QueryInput {
		TableName: aws.String(ddbTablename),
		KeyConditionExpression: aws.String("(#MealTime = :mor OR #MealTime = :lun OR #MealTime = :din) AND #Date >= :fromdate"),
		ExpressionAttributeNames: map[string]*string {
			"#MealTime": aws.String("MealTime"),
			"#Date": aws.String("Date"),
			"#Day_of_week": aws.String("Day_of_week"),
			"#MealMethod": aws.String("MealMethod"),		
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue {
			":mor": {
				S: aws.String("MOR"),
			},
			":lun": {
				S: aws.String("LUN"),
			},
			":din": {
				S: aws.String("DIN"),
			},
			":fromdate": {
				N: aws.String("20170901"),
			},
		},
		ProjectionExpression: aws.String("#MealTime, #Date, #Day_of_week, #MealMethod"),
		ScanIndexForward: aws.Bool(false),

	}

	queryItem, queryErr := svc.Query(queryParams)
	if queryErr != nil {
		panic(queryErr)
	}

	// convert json
//	mhdatas := make([]*MHdata, 0)
	if err := dynamodbattribute.UnmarshalListOfMaps(queryItem.Items, &mhdatas); err != nil {
		log.Println("[Unmarshal Error]", err)
		panic(queryErr)
	}

	return fmt.Sprintln(queryItem.Items), nil
}

func Handler() ([]MHdata, error) {
	// 処理部分

	// 変数初期化
	data := make([]MHdata,0)
	msg, err := Scanmealhistory(&data)

	// output data
//	var output []string
//	for _, i := range data {
//		append(output,string(data[i]))
//	}

	log.Println(msg)

	return data, err
}

func main() {
    lambda.Start(Handler)
}