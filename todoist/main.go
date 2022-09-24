package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaResponse events.APIGatewayProxyResponse

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (LambdaResponse, error) {
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("EVENT: %s", eventJson)
	url := "https://api.todoist.com/rest/v2/tasks"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", os.Getenv("TODOIST_AUTH"))

	todoistResp, _ := client.Do(req)
	//We Read the response body on the line below.
	todoistBody, err := ioutil.ReadAll(todoistResp.Body)
	if err != nil {
		log.Fatalln("Couldn't read body", err.Error())
	}
	// Try to convert the body to Todoist Task
	var tasks []TodoistTask
	err = json.Unmarshal(todoistBody, &tasks)
	if err != nil {
		log.Fatalln("Couldn't Unmarshal body", err.Error())
	}

	trackeableTasks := filterTrackeableTasks(tasks)
	// for _, task := range trackeableTasks {
	// 	log.Println(task.Content, task.wasDoneToday())
	// }

	// TODO create the entrypoint function (the lambda handler) that will check at the end of the day
	//  if the task "wasDoneToday" update the streaks table, and log the entry on the taskhistory table.
	// on the streal table, save last streak, last streak start, last streak end. Because sometimes I forget about
	// checking todoist at the end of the day and then I update the next day and I shouldn't lose the streak in these cases.
	// As said in the beginning of this doc, the ideia is to create a lambda function that will run everyday at 23:55
	// and update the tables
	body, err := json.Marshal(trackeableTasks)
	if err != nil {
		log.Fatalln("Couldn't marshal trackeable tasks", err.Error())
	}
	lambdaResp := LambdaResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
	return lambdaResp, nil
}
func main() {
	lambda.Start(Handler)
}
