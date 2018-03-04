package serverless

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/khezen/espipe/configuration"
	"github.com/khezen/espipe/document"
	"github.com/khezen/espipe/errors"
	"github.com/khezen/espipe/indexer"
	"github.com/khezen/espipe/template"
)

// IndexHandler index doument in bulk request to elastic search
func IndexHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	configuration.Set("etc/config.json")
	config, err := configuration.Get()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: errors.HTTPStatusCode(err)}, err
	}
	if !config.Redis.Enabled {
		err := fmt.Errorf("serverless mode require Redis to be enabled")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: errors.HTTPStatusCode(err)}, err
	}
	config.Redis.AutoFlush = false
	indexer, err := indexer.New(*config)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: errors.HTTPStatusCode(err)}, err
	}
	urlSplit := strings.Split(strings.Trim(strings.ToLower(request.Path), "/"), "/")
	if len(urlSplit) != 3 {
		return events.APIGatewayProxyResponse{Body: errors.ErrPathNotFound.Error(), StatusCode: errors.HTTPStatusCode(errors.ErrPathNotFound)}, err
	}
	docTemplate := template.Name(urlSplit[1])
	docType := document.Type(urlSplit[2])
	docByte := []byte(request.Body)
	err = indexer.Index(docTemplate, docType, docByte)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: errors.HTTPStatusCode(err)}, err
	}
	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func main() {
	lambda.Start(IndexHandler)
}
