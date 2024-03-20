package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DmitriyKomarovCoder/flood-control/flood"
	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	f := flood.NewFloodControl(1, 1*time.Second, client)
	f.Check(context.Background(), 1)
	i := int64(1)
	for {
		i++
		fmt.Println(f.Check(context.Background(), 1))
		fmt.Println(f.Check(context.Background(), 1))
		fmt.Println(f.Check(context.Background(), 1))
		time.Sleep(1 * time.Second)
	}
}
