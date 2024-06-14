package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/xissg/online-judge/internal/model/entity"
	"strconv"
)

type RESPONSE interface {
	entity.Question | entity.Submit
}

func cacheOrUpdateData[T RESPONSE](rdb *RedisClient, key string, scores []int, ids []int, ts []*T) error {
	pipe := rdb.client.Pipeline()
	for i, id := range ids {
		filed := strconv.Itoa(id)
		member := redis.Z{
			Score:  float64(scores[i]),
			Member: filed,
		}
		_, err := pipe.ZAdd(key, member).Result()
		if err != nil {
			return err
		}
		content := fmt.Sprintf("%v_content", key)
		data, err := json.Marshal(ts[i])
		_, err = pipe.HSet(content, filed, data).Result()
		if err != nil {
			return err
		}
	}
	_, err := pipe.Exec()
	return err
}

func getDataList[T RESPONSE](rdb *RedisClient, key string, page, pageSize int) (ts []*T) {
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	ids, err := rdb.client.ZRange(key, int64(start), int64(end)).Result()
	if err != nil {
		return nil
	}

	var t T
	content := fmt.Sprintf("%v_content", key)
	for i := range ids {
		data, _ := rdb.client.HGet(content, ids[i]).Result()
		err := json.Unmarshal([]byte(data), &t)
		if err != nil {
			return nil
		}

		temp := t
		ts = append(ts, &temp)
	}

	return ts
}

func deleteDataById(redis *RedisClient, key string, id int) error {
	field := strconv.Itoa(id)
	_, err := redis.client.ZRem(key, field).Result()
	if err != nil {
		return err
	}

	content := fmt.Sprintf("%v_content", key)
	_, err = redis.client.HDel(content, field).Result()
	if err != nil {
		return err
	}
	return nil
}
