package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/satori/go.uuid"
	"github.com/serverless-crud-go/src/package/database"
	"log"
)

type Todo struct {
	ID          string  `json:"id"`
	Description string 	`json:"description"`
	Done        bool   	`json:"done"`
	CreatedAt   string 	`json:"created_at"`
}


// TODO - move code dealing with Database interaction to repository

func AddTodo(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ddb, err := database.ConnectionToDatabase(os.Getenv("AWS_REGION"))
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("AddTodo")

	var (
		id = uuid.Must(uuid.NewV4(), nil).String()
		tableName = aws.String(os.Getenv("TODOS_TABLE_NAME"))
	)

	// Initialize todo
	todo := &Todo{
		ID:					id,
		Done:				false,
		CreatedAt:			time.Now().String(),
	}

	// Parse request body
	json.Unmarshal([]byte(request.Body), todo)

	// Write to DynamoDB
	item, _ := dynamodbattribute.MarshalMap(todo)
	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: tableName,
	}
	if _, err := ddb.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{ // Error HTTP response
			Body: err.Error(),
			StatusCode: 500,
		}, nil
	} else {
		body, _ := json.Marshal(todo)
		return events.APIGatewayProxyResponse{ // Success HTTP response
			Body: string(body),
			StatusCode: 200,
		}, nil
	}
}

func main() {

	lambda.Start(AddTodo)
}