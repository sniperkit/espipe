package serverless

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/khezen/espipe/configuration"
	"github.com/khezen/espipe/dispatcher"
)

// FlushHandler index doument in bulk request to elastic search
func FlushHandler() (string, error) {
	config, err := configuration.Get()
	if err != nil || !config.Redis.Enabled {
		configuration.Set("etc/config.json")
		config, err = configuration.Get()
		if err != nil {
			return err.Error(), err
		}
	}
	if !config.Redis.Enabled {
		err := fmt.Errorf("serverless mode require Redis to be enabled")
		return err.Error(), err
	}
	config.Redis.AutoFlush = false
	buffer := dispatcher.RedisBuffer(nil, config)
	err = buffer.Flush()
	if err != nil {
		return err.Error(), err
	}
	return "", nil
}

func main() {
	lambda.Start(FlushHandler)
}
