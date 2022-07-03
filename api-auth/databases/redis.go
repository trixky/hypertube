package databases

import (
	"strconv"

	"github.com/go-redis/redis"
)

const (
	REDIS_SEPARATOR     = "."
	REDIS_PATTERN_token = "token" + REDIS_SEPARATOR
)

func AddToken(user_id int64, token string) error {
	err := DBs.Redis.Append(REDIS_PATTERN_token+strconv.Itoa(int(user_id)), token+REDIS_SEPARATOR).Err()

	return err
}

func InitRedis() error {
	DBs.Redis = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if _, err := DBs.Redis.Ping().Result(); err != nil {
		return err
	}

	return nil
}
