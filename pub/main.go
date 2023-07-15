package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var count = 0
var countChan = make(chan int)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	go func() {
		for {
			count := <-countChan
			err := rdb.Publish(ctx, "mychannel", count).Err()
			if err != nil {
				panic(err)
			}
		}
	}()

	http.HandleFunc("/", countMiddleware(handleGet))
	http.ListenAndServe(":4000", nil)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("hello!")
}

func countMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count += 1
		countValue := count

		go func() {
			countChan <- countValue
		}()

		f(w, r)
	}
}
