package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
)

func newItemsKey() string {
	return "cache:newItems:ids"
}

type redisGetNewItemsIDsParam struct {
	LastID        int64
	LastCreatedAt int64
	Size          int64
}

func addNewItems(id, createdAt int64) {
	redisClient.ZAdd(newItemsKey(), redis.Z{
		Score:  float64(createdAt),
		Member: id,
	})
}

func delNewItems(id int64) {
	redisClient.ZRem(newItemsKey(), id)
}

func getNewItemsIDsFromRedis(param redisGetNewItemsIDsParam) []int64 {
	if param.LastID > 0 && param.LastCreatedAt > 0 {
		zs, err := redisClient.ZRevRangeByScoreWithScores(newItemsKey(), redis.ZRangeBy{
			Max:    fmt.Sprintf("%d", param.LastCreatedAt),
			Min:    "0",
			Offset: 0,
			Count:  param.Size + 100,
		}).Result()
		if err != nil {
			log.Print(err)
			return nil
		}

		align := false
		ids := make([]int64, 0, len(zs))
		for _, z := range zs {
			sid := z.Member.(string)
			id, _ := strconv.ParseInt(sid, 10, 64)

			if !align && param.LastCreatedAt == int64(z.Score) && id == param.LastID {
				align = true
				continue
			}

			ids = append(ids, id)
			if int64(len(ids)) == param.Size+1 {
				break
			}
		}
		return ids
	}

	sids, err := redisClient.ZRevRange(newItemsKey(), 0, param.Size+1).Result()
	if err != nil {
		log.Print(err)
		return nil
	}

	ids := make([]int64, 0, len(sids))
	for _, sid := range sids {
		id, _ := strconv.ParseInt(sid, 10, 64)
		ids = append(ids, id)
	}

	return ids
}

//type redisGetNewCategoryItemsIDsParam struct {
//	CategoryID    int
//	LastID        int
//	LastCreatedAt int
//	Size          int
//}
//
//func getNewCategoryItemIDsFromRedis() []uint64 {
//
//}
