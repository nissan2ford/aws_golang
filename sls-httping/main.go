package main

import (
	"fmt"
	"httpfunc"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	ReqURL   string        `json:"requrl"`
	Response string        `json:"response"`
	Duration time.Duration `json:"duration"`
	DoTime   int           `json:"dotime"`
}

func Handler() (Response, error) {

	urls := []string{
		//		"http://192.168.11.1",
		"http://www.intellilink.co.jp",
		//		"http://github.com",
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

	return Response{
		ReqURL:   requrlstr,
		Response: respstr,
		Duration: duratime,
		DoTime:   dotimeint,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
