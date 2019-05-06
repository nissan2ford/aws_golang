package main

import (
	"fmt"
	"httpfunc"
	"time"
)

func main() {
	urls := []string{
		"http://192.168.11.1",
		"http://www.yahoo.co.jp",
		"http://github.com",
	}

	// make channel
	responseChan := make(chan string)
	durationChan := make(chan time.Duration)
	reqUrlChan := make(chan string)

	for _, url := range urls {

		// go routine
		go func(url string) {

			requrl,rescode,duration := httpfunc.ConnHttp(url)

			reqUrlChan <- requrl
			responseChan <- rescode
			durationChan <- duration
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Println(<-reqUrlChan,<-responseChan,"TAT=",<-durationChan)
	}
}
