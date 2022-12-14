package main

import (
	"encoding/json"
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
	// Make the call to the DAO with params found in the path
	// fmt.Println("Path vars: ", request.PathParameters["year"], " ", parseSlug(request.PathParameters["title"]))
	// item, err := GetByYearTitle(request.PathParameters["year"], parseSlug(request.PathParameters["title"]))
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to find Item, %v", err))
	// }

	// // Make sure the Item isn't empty
	// if item.Year <= 0 {
	// 	fmt.Println("Could not find movie")
	// 	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, nil
	// }

	// // Log and return result
	// jsonItem, _ := json.Marshal(item)
	// stringItem := string(jsonItem) + "\n"
	// fmt.Println("Found item: ", stringItem)
	allTasks := GetAll()
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
