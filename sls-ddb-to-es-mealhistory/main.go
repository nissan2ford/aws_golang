package main

import (
	"fmt"

	"log"

	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"


	"net/http"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"gopkg.in/olivere/elastic.v3"
	"github.com/edoardo849/apex-aws-signer"
)

// define elasticsearch setting
const (
	// dynamoDB Table Name
	ddbTablename = "meal_history"

	// elasticsearch setting
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

func Scanmealhistory(mhdatas *[]MHdata)(error) {
	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess,aws.NewConfig().WithRegion("ap-northeast-1"))

	// Scan Parameter
	scanParams := &dynamodb.ScanInput {
		TableName: aws.String(ddbTablename),
	}

	scanItem, err := svc.Scan(scanParams)

	if err != nil {
		panic(err)
	}

	// convert struct
	err = dynamodbattribute.UnmarshalListOfMaps(scanItem.Items, &mhdatas)
	if err != nil {
		 log.Printf("failed to unmarshal Query result items, %v", err)
	}

	return nil
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

	return string(b), nil
}

func Handler() (string, error) {
	// 処理部分

	// 変数初期化
	data := make([]MHdata,0)

	// Scan from DynamoDB
	err := Scanmealhistory(&data)

	// print mhdatas
	log.Println(data)

	// Insert elasticsearch from mhdatas
	for _, mhdata := range mhdatas {

		msg, err := PutES(mhdata)

		if err != nil {
			panic(err)
		}
	}

	// data count
	count := range mhdatas

	return fmt.Printf("%d mealhistory datas was inserted",count), err
}

func main() {
    lambda.Start(Handler)
}