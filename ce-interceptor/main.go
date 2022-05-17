package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// Start the CloudEvent handler, goal is to take the db CUD change events,
// transform and apply the changes to redis
func receiver(ev cloudevents.Event) {
	// Extract the cloudevent body
	body := make(map[string]interface{})
	_ = ev.DataAs(&body)

	// Retrieve the payload
	payload := body["payload"].(map[string]interface{})
    var payloadBefore map[string]interface{}
    var payloadAfter map[string]interface{}

    if payload["before"] != nil {
        payloadBefore = payload["before"].(map[string]interface{})
    }

    if payload["after"] != nil {
        payloadAfter = payload["after"].(map[string]interface{})
    }

    // This is an update or add operation. For our example, we don't care about
    // comparing the the prior -> future state so we will only apply the future state.
    if payloadAfter != nil {
        name := fmt.Sprintf("%v", payloadAfter["name"])
        log.Printf("setting: %s [%d]", name, payloadAfter["vote"])
        err := redisClient.Set(context.Background(), name, payloadAfter["vote"], 0).Err()
        if err != nil {
            log.Printf("Err: unable to add/update: %v", err)
        }
    }

    // This is a delete operation
    if payloadAfter == nil && payloadBefore != nil {
        name := fmt.Sprintf("%v", payloadBefore["name"])
        log.Printf("deleting: %s", name)
        err := redisClient.Del(context.Background(), name).Err()
        if err != nil {
            log.Printf("Err: unable to delete: %v", err)
        }
    }
}

// Start the CloudEvent server, and register the connection with Redis
func main() {
	log.Println("Starting Change Event Capture...")
	// connect to Redis as passed in by the environment variable REDIS_HOST
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // ignore password/clustering for this example
		DB:       0,  // ignore custom db for now
	})

	// connect to the cloudevent server
	ce, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("Unable to create server: %v", err)
	}

	log.Fatal(ce.StartReceiver(context.Background(), receiver))
}
