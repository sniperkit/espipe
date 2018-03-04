build:
	go get github.com/aws/aws-sdk-go/aws/signer/v4
	go get github.com/aws/aws-sdk-go/aws/credentials
	go get github.com/google/uuid
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-lambda-go/events
	go get github.com/go-redis/redis
	env GOOS=linux go build -ldflags="-s -w" -o bin/index index.go