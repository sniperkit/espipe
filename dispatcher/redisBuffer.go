package dispatcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/khezen/espipe/configuration"
	"github.com/khezen/espipe/document"
	"github.com/khezen/espipe/elastic"
	"github.com/khezen/espipe/template"
)

type redisBuffer struct {
	redis     *redis.Client
	Template  *template.Template
	elastic   *elastic.Client
	Kill      chan error
	documents []document.Document
	sizeKB    float64
	mutex     sync.RWMutex
}

// RedisBuffer -
func RedisBuffer(template *template.Template, config *configuration.Configuration) Buffer {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.Partition,
	})
	elasticClient := elastic.NewClient(config)
	buffer := &redisBuffer{
		redisClient,
		template,
		elasticClient,
		make(chan error),
		make([]document.Document, 0),
		0,
		sync.RWMutex{},
	}
	return buffer
}

func (b *redisBuffer) Append(msg document.Document) error {
	key := fmt.Sprintf("espipe:%s.%s", msg.Template.Name, msg.Type)
	request, err := msg.Request()
	if err != nil {
		return err
	}
	res := b.redis.Append(key, string(request))
	err = res.Err()
	if err != nil {
		return err
	}
	return nil
}

func (b *redisBuffer) Flush() error {
	res := b.redis.Keys("^espipe:*")
	keys, err := res.Result()
	if err != nil {
		return err
	}
	bulk := make([]byte, 0, int(b.sizeKB)+len(b.documents)*150)
	for _, key := range keys {
		res := b.redis.Get(key)
		if res.Err() != nil {
			return res.Err()
		}
		requests := res.String()
		bulk = append(bulk, []byte(requests)...)
	}
	err = b.elastic.Bulk(bulk)
	if err != nil {
		return err
	}
	for _, key := range keys {
		delRes := b.redis.Del(key)
		err = delRes.Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// Flusher flushes every {configuration/config.go::Template.TimerMS}
func (b *redisBuffer) Flusher() func() {
	ticker := time.NewTicker(time.Duration(b.Template.FlushPeriodMS) * time.Millisecond)
	return func() {
		for {
			select {
			case <-b.Kill:
				return
			case <-ticker.C:
				go func() {
					err := b.Flush()
					if err != nil {
						fmt.Println(err)
					}
				}()
				break
			}
		}
	}
}
