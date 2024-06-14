package redis

func (redis *RedisClient) HIncrBy(key string, field string, incr int64) error {
	_, err := redis.client.HIncrBy(key, field, incr).Result()
	return err
}

func (redis *RedisClient) HGetAll(key string) (map[string]string, error) {
	data, err := redis.client.HGetAll(key).Result()
	return data, err
}
