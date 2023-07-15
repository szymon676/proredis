package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var count = 0

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pubsub := rdb.Subscribe(ctx, "mychannel")
	defer pubsub.Close()

	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			convmsg, _ := strconv.Atoi(msg.Payload)
			count = convmsg
			fmt.Println(count)
		}
	}()

	http.HandleFunc("/count", handleCount)
	http.ListenAndServe(":3000", nil)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&count)
}
