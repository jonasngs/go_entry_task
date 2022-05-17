package services

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	pb "github.com/jonasngs/go_entry_task/grpc"
	"github.com/jonasngs/go_entry_task/tcpserver/storage"
)

type CacheInterface interface {
	checkCache(sessionToken string) (*pb.User, error)
	updateCache(sessionToken string, user *pb.User) error
	deleteFromCache(sessionToken string)
}

type CacheService struct {
	cache storage.RedisCache
}

var ctx = context.Background()

func InitializeCacheService(cache storage.RedisCache) CacheInterface {
	return CacheService{cache: cache}
}

func (cs CacheService) checkCache(sessionToken string) (*pb.User, error) {
	cachedUser, err := cs.cache.RDB.Get(ctx, sessionToken).Result()
	if err != nil {
		log.Printf("Unable to verify cache: %s\n", err)
		return nil, err
	}
	user := &pb.User{}
	err = proto.Unmarshal([]byte(cachedUser), user)
	if err != nil {
		log.Printf("Unable to unmarshal user: %s\n", err)
		return nil, err
	}
	return user, nil
}

func (cs CacheService) updateCache(sessionToken string, user *pb.User) error {

	serializedUser, err := proto.Marshal(user)
	if err != nil {
		log.Printf("Unable to marshal user: %s\n", err)
		return err
	}

	err = cs.cache.RDB.Set(ctx, sessionToken, serializedUser, time.Hour*5).Err()
	if err != nil {
		log.Printf("Unable to cache session: %s\n", err)
		return err
	}
	return nil
}

func (cs CacheService) deleteFromCache(sessionToken string) {
	cs.cache.RDB.Del(ctx, sessionToken)
}
