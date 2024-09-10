package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = map[string]Book{
	"1": {ID: "1", Title: "Song of Achilles", Author: "Madeline Miller"},
	"2": {ID: "2", Title: "The Other Boleyn Girl", Author: "Philippa Gregory"},
}

func extractBookID(path string) (string, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 2 {
		return parts[1], true
	}
	return "", false
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch event.HTTPMethod {
	case http.MethodGet:
		if id, ok := extractBookID(event.Path); ok {
			return getBook(id)
		}
		return listBooks()
		// case http.MethodPost:
		// 	return createBook(event)
		// case http.MethodPut:
		// 	if id, ok := event.PathParameters["id"]; ok {
		// 		return updateBook(id, event)
		// 	}
		// case http.MethodDelete:
		// 	if id, ok := event.PathParameters["id"]; ok {
		// 		return deleteBook(id)
		// 	}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       "Method Not Allowed",
	}, nil
}

func main() {
	lambda.Start(handler)
}

func listBooks() (events.APIGatewayProxyResponse, error) {
	bookList := []Book{}
	for _, book := range books {
		bookList = append(bookList, book)
	}

	body, err := json.Marshal(bookList)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling books",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func getBook(id string) (events.APIGatewayProxyResponse, error) {
	if book, ok := books[id]; ok {
		body, err := json.Marshal(book)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Error marshalling book",
			}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       "Book not found",
	}, nil
}
