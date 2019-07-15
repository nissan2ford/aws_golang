package main

import (
	"fmt"
	"httpfunc"
	"strconv"
	"time"

	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type dynamoData struct {
	ReqURL   string `dynamo:"Req_URL"`
	Response string `dynamo:"Res_Code"`
	Duration int64  `dynamo:"Duration"`
	DoTime   int    `dynamo:"Do_Time"`
}

// define dynamoDB setting
const (
	// dynamoDB Table Name
	hcTablename = "IL_hc"

	// AWS Region
	region = "ap-northeast-1"
)

func Handler() (dynamoData, error) {

	urls := []string{
		os.Getenv("RequestURL"),
	}

	// make channel
	responseChan := make(chan string)
	durationChan := make(chan time.Duration)
	reqUrlChan := make(chan string)

	// dotime var
	dotime := time.Now()

	for _, url := range urls {

		// go routine
		go func(url string) {

			requrl, rescode, duration := httpfunc.ConnHttp(url)

			reqUrlChan <- requrl
			responseChan <- rescode
			durationChan <- duration

		}(url)
	}

	// ruterun value
	var requrlstr string
	var respstr string
	var duratime time.Duration

	for i := 0; i < len(urls); i++ {
		requrlstr = <-reqUrlChan
		respstr = <-responseChan
		duratime = <-durationChan
		fmt.Println(requrlstr, respstr, "TAT=", duratime, "DoTime=", dotime)
	}

	// change format time var
	const layout = "20060102150405"
	dotimeint, _ := strconv.Atoi(dotime.Format(layout))

	duratimeint := int64(duratime)

	// set dynamoData
	putdata := dynamoData{
		ReqURL:   requrlstr,
		Response: respstr,
		Duration: duratimeint,
		DoTime:   dotimeint,
	}

	// session
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess, aws.NewConfig().WithRegion(region))

	// パラメータ
	params := &dynamodb.PutItemInput{
		TableName: aws.String(hcTablename), // table name
		Item: map[string]*dynamodb.AttributeValue{ // Input Datas
			"Req_URL": {
				S: aws.String(putdata.ReqURL),
			},
			"Res_Code": {
				S: aws.String(putdata.Response),
			},
			"Duration": {
				N: aws.String(strconv.FormatInt(putdata.Duration, 10)), // int64 -> string
			},
			"Do_Time": {
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
