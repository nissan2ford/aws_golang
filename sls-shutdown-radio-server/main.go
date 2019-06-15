package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"os"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/awserr"

)

type Response struct {
	Message string `json:"message"`
}

func Handler() (Response, error) {

	insidstr := os.Getenv("InstanceID")

	svc := ec2.New(session.New())

	actionstr := os.Getenv("Action")

	if actionstr == "start" {
		
		input := &ec2.StartInstancesInput{
	    InstanceIds: []*string{
  	      aws.String(insidstr),
    	},
		}

		result, err := svc.StartInstances(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			
			fmt.Println(result)
		}

	} else if actionstr == "stop" {
	
		input := &ec2.StopInstancesInput{
	    InstanceIds: []*string{
  	      aws.String(insidstr),
    	},
		}

		result, err := svc.StopInstances(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}

			fmt.Println(result)

		}

	}

	returnmsg := actionstr + " instance successfully!"

	return Response{
		Message: returnmsg,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
