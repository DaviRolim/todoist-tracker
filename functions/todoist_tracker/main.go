package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	invoke "github.com/aws/aws-sdk-go/service/lambda"
)

type Payload struct {
	// You can also include more objects in the structure like below,
	// but for my purposes body was all that was required
	// Method string `json:"httpMethod"`
	Body string `json:"body"`
}

const TODOIST_URL = "https://api.todoist.com/rest/v2/tasks"

var sess = session.Must(session.NewSession())
var svc = invoke.New(sess)

type LambdaResponse events.APIGatewayProxyResponse

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (LambdaResponse, error) {
	log.Printf("EVENT: %v", event)
	// Call the todoist public API to get my tasks
	tasks := getTasksFromTodoist()

	trackeableTasks := filterTrackeableTasks(tasks)

	body, err := json.Marshal(trackeableTasks)
	if err != nil {
		log.Fatalln("Couldn't marshal trackeable tasks", err.Error())
	}
	// Send trackeable tasks to a lambda that has connection to a DB (dynamo in this case) to save the tasks
	sendTasksToLambdaToPersist(trackeableTasks)

	// Send a response with all trackeable tasks
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

func getTasksFromTodoist() []TodoistTask {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", TODOIST_URL, nil)
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

	return tasks
}

func main() {
	lambda.Start(Handler)
}
func sendTasksToLambdaToPersist(trackeableTasks []TodoistTask) {
	for _, task := range trackeableTasks {
		if !task.wasDoneToday() {
			continue
		}
		payload := buildPayloadFromTask(task)
		input := &invoke.InvokeInput{
			FunctionName:   aws.String("addCompletedTask"),
			InvocationType: aws.String("Event"),
			Payload:        payload,
		}
		_, err := svc.Invoke(input)
		if err != nil {
			fmt.Println("error Invoking another function")
			fmt.Println(err.Error())
		}
	}
}

func buildPayloadFromTask(task TodoistTask) []byte {
	const layout = "2006-01-02"
	taksDueTime, err := time.Parse(layout, task.Due.Date)
	// if the task was done today, the next DueDate will be tomorrow, so I'll decrease one day from due date to get the "done" date.
	// time.Now().Format("2006-01-02") should also work.
	//today := taksDueTime.Add(time.Hour * 24 * -1)
	today := taksDueTime.AddDate(0, 0, -1)
	var taskJson = make(map[string]any)
	taskJson["name"] = task.Content
	taskJson["date"] = today.Format("2006-01-02")
	if task.Description != "" {
		addMeasurementsToTaskJsonFromTaskDescription(taskJson, task)
	}
	body, err := json.Marshal(taskJson)
	if err != nil {
		log.Fatalln("Error parsing body", err.Error())
	}
	p := Payload{
		Body: string(body),
	}
	payload, err := json.Marshal(p)
	if err != nil {
		log.Fatalln("Error parsing payload", err.Error())
	}
	return payload
}

func addMeasurementsToTaskJsonFromTaskDescription(taskJson map[string]any, task TodoistTask) {
	var description map[string]any
	err := json.Unmarshal([]byte(task.Description), &description)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(description)
	taskJson["measurementType"] = description["measurementType"]
	taskJson["measurementValue"] = description["measurementValue"]
}

func filterTrackeableTasks(tasks []TodoistTask) []TodoistTask {
	var result []TodoistTask
	for _, task := range tasks {
		if task.isTrackeable() {
			result = append(result, task)
		}
	}
	return result
}
