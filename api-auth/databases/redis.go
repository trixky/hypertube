package databases

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

const (
	REDIS_EXTERNAL_google   = "google"
	REDIS_SEPARATOR         = "."
	REDIS_PATTERN_KEY_token = "token"
)

type Token_info struct {
	Id       int64
	External string
}

func AddToken(user_id int64, token string, external string) error {
	if len(external) == 0 {
		external = EXTERNAL_none
	}

	err := DBs.Redis.Append(REDIS_PATTERN_KEY_token+REDIS_SEPARATOR+token+REDIS_SEPARATOR+strconv.Itoa(int(user_id)), external).Err()

	return err
}

func RetrieveToken(token string) (*Token_info, error) {
	keys, err := DBs.Redis.Keys(REDIS_PATTERN_KEY_token + REDIS_SEPARATOR + token + REDIS_SEPARATOR + "*").Result()

	if err != nil {
		return nil, err
	}

	key_nbr := len(keys)

	if key_nbr != 1 {
		return nil, fmt.Errorf("expected one session for one token (%d finded)", key_nbr)
	}

	id, err := strconv.Atoi(strings.SplitN(keys[0], ".", 3)[2])

	if err != nil {
		return nil, fmt.Errorf("token backup is broken")
	}

	external, err := DBs.Redis.Get(keys[0]).Result()

	if err != nil {
		return nil, fmt.Errorf("token extraction failed")
	}

	return &Token_info{
		Id:       int64(id),
		External: external,
	}, nil
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
