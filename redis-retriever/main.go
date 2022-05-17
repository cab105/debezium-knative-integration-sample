package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func redisHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all projects and votes to render a plain-text response
	w.Header().Set("content-type", "text/plain")

	projects := make(map[string]string, 0)
	ctx := context.Background()
	iter := redisClient.Scan(ctx, 0, "*", 0).Iterator()

	for iter.Next(ctx) {
		k := iter.Val()
		vote, err := redisClient.Get(ctx, k).Result()
		if err != nil {
			log.Panic("err retrieving value: %v", err)
		}
		projects[k] = vote
	}

	if err := iter.Err(); err != nil {
		log.Panic("err retrieving keys: %v", err)
	}

	var buf string
	for k, v := range projects {
		buf = buf + "\n" + fmt.Sprintf("%25s | %s\n", k, v)
	}

	_, err := w.Write([]byte(buf))
	if err != nil {
		log.Panic("err writing output: %v", err)
	}
}

func main() {
	// connect to Redis as passed in by the environment variable REDIS_HOST
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Panic("REDIS_HOST must be defined")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // ignore password/clustering for this example
		DB:       0,  // ignore custom db for now
	})


	// Start the web server
	http.HandleFunc("/", redisHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
