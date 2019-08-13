package main

import (
	"fmt"
	"strconv"
//	"time"
	"os"
	
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// define input data type
type inputData struct{
 	Health	string `json:"health"`
 	Perce	int	`json:"percentage"`
 	Plugged	string	`json:"plugged"`
 	Status	string	`json:"status"`
 	Tempe	float64	`json:"temperature"`
 	DoTime	int	`json:"dotime"`
 }

// define dynamoDB data type
type dynamoData struct {
 	Health	string `dynamo:"health"`
 	Perce	int	`dynamo:"percentage"`
 	Plugged	string	`dymamo:"plugged"`
 	Status	string	`dynamo:"status"`
 	Tempe	float64	`dynamo:"temperature"`
	DoTime   int    `dynamo:"dotime"`
}

// define dynamoDB setting
const (
	// AWS Region
	region = "ap-northeast-1"
)

func Handler(indata inputData) (dynamoData, error) {

	// dynamoDB Table Name
	ddbTablename := os.Getenv("dDBTablename")

	// dotime var
//	dotime := time.Now()

	// change format time var
//	const layout = "20060102150405"
//	dotimeint, _ := strconv.Atoi(dotime.Format(layout))

	// set dynamoData
	putdata := dynamoData{
		Health:	indata.Health,
		Perce:		indata.Perce,
		Plugged:	indata.Plugged,
		Status:		indata.Status,
		Tempe:	indata.Tempe,
		DoTime:	indata.DoTime,
//		DoTime:	dotimeint,
	}

	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess, aws.NewConfig().WithRegion(region))

	// パラメータ
	params := &dynamodb.PutItemInput{
		TableName: aws.String(ddbTablename), // table name
		Item: map[string]*dynamodb.AttributeValue{ // Input Datas
			"health": {
				S: aws.String(putdata.Health),
			},
			"percentage": {
				N: aws.String(strconv.Itoa(putdata.Perce)),
			},
			"plugged":{
				S: aws.String(putdata.Plugged),
			},
			"status":{
				S: aws.String(putdata.Status),
			},
			"temperature": {
				N: aws.String(strconv.FormatFloat(putdata.Tempe, 'e',14,64)), // float64 -> string
			},
			"dotime": {
				N: aws.String(strconv.Itoa(putdata.DoTime)),
			},
		},
	}

	// PutItemの実行
	resp, err := svc.PutItem(params)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	return putdata, nil
}

func main() {
	lambda.Start(Handler)
}
