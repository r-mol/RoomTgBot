package main

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v9"
)

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()

	if err != nil {
		panic(err)
	}

	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()

	switch err {
	case redis.Nil:
		fmt.Println("key2 does not exist")
	case nil:
		panic(err)
	default:
		fmt.Println("key2", val2)
	}

	// Output: key value
	// key2 does not exist
	fmt.Println("Redis test success")
}
