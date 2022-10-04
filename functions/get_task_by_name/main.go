package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

// Parse slug into a space separated string
func parseSlug(orig string) (retval string) {
	retval = strings.Replace(orig, "-", " ", -1)
	retval = strings.Replace(retval, "+", " ", -1)
	return retval
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("Path vars: ", request.QueryStringParameters["name"])
	allTasks := GetAllByName(request.QueryStringParameters["name"])
	body, err := json.Marshal(allTasks)
	if err != nil {
		log.Fatal("Unable to parse body")
	}

	return events.APIGatewayProxyResponse{
		Body: string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
