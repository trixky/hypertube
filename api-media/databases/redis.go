package databases

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

const (
	REDIS_SEPARATOR          = "."
	REDIS_PATTERN_KEY_token  = "token"
	REDIS_PATTERN_KEY_search = "search"
)

type Token_info struct {
	Id       int64
	External string
}

func RetrieveToken(token string) (*Token_info, error) {
	keys, err := DBs.Redis.Keys(REDIS_PATTERN_KEY_token + REDIS_SEPARATOR + token + REDIS_SEPARATOR + "*").Result()

	if err != nil {
		return nil, err
	}

	key_nbr := len(keys)

	if key_nbr != 1 {
		return nil, fmt.Errorf("expected one session for one token (%d found)", key_nbr)
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

func AddSearch(path *string, results *string) error {
	// Save in redis
	err := DBs.Redis.Append(REDIS_PATTERN_KEY_search+*path, *results).Err()
	return err
}

func RetrieveSearch(path *string) (string, error) {
	// Check if the search exists on redis
	results, err := DBs.Redis.Get(REDIS_PATTERN_KEY_search + *path).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", fmt.Errorf("search extraction failed")
		}
	}

	// Return the results if they exists
	// -- Convert in the caller, to avoid clutter here
	return results, nil
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
