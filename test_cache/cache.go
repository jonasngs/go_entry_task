package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	pb "github.com/jonasngs/go_entry_task/grpc"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// err = rdb.Set(ctx, "LOL", GenerateSessionToken(), 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// rdb.Del(ctx, "key")

	// val, err := rdb.Get(ctx, "LOL").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("LOL", val)

	// val, err = rdb.Get(ctx, "key").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key does not exist")
	// }
	// fmt.Println("key", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
	// Output: key value
	// key2 does not exist
	updateCache(rdb)

}

func updateCache(rdb *redis.Client) error {
	user := &pb.User{
		UserId:         1,
		Username:       "user",
		Password:       "password",
		Nickname:       "nickname",
		ProfilePicture: "pp",
	}
	serializedUser, err := proto.Marshal(user)
	if err != nil {
		log.Printf("Unable to marshal user: %s\n", err)
		return err
	}
	token := GenerateSessionToken()
	err = rdb.Set(ctx, token, serializedUser, 1000000000000).Err()
	if err != nil {
		log.Printf("Unable to cache session: %s\n", err)
		return err
	}

	// val2, err := rdb.Get(ctx, token).Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println(token)
	// 	fmt.Println(val2)
	// 	fmt.Println()
	// }

	cachedUser, err := rdb.Get(ctx, token).Result()
	if err != nil {
		fmt.Print("LOL")
	}

	user1 := &pb.User{}
	err = proto.Unmarshal([]byte(cachedUser), user1)
	fmt.Println(user1)

	// err = rdb.Set(ctx, "LOL", GenerateSessionToken(), 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	newUser := user1
	newUser.Nickname = "LLLLLLLLLOOOOOLLLL"
	finalUser, err := proto.Marshal(newUser)
	if err != nil {
		log.Printf("Unable to marshal user: %s\n", err)
		return err
	}
	err = rdb.Set(ctx, token, finalUser, 1000000000000).Err()
	if err != nil {
		log.Printf("Unable to cache session: %s\n", err)
		return err
	}
	return nil
}

func GenerateSessionToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Unable to generate random byte sequence: %s\n", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
