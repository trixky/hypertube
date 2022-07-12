package databases

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	REDIS_EXTERNAL_google            = "google"
	REDIS_SEPARATOR                  = "."
	REDIS_PATTERN_KEY_token          = "token"
	REDIS_PATTERN_KEY_password_token = "password_token"
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

func AddPasswordToken(user_id int64, token string) error {
	err := DBs.Redis.Set(REDIS_PATTERN_KEY_password_token+REDIS_SEPARATOR+token+REDIS_SEPARATOR+strconv.Itoa(int(user_id)), "", 10*time.Minute).Err()

	return err
}

func RetrievePasswordToken(token string, delete bool) (*Token_info, error) {
	user_id, _, _, err := retrievePattern(REDIS_PATTERN_KEY_password_token, token)

	if err != nil {
		return nil, err
	}

	if delete {
		if err := DBs.Redis.Del(REDIS_PATTERN_KEY_password_token + REDIS_SEPARATOR + token + REDIS_SEPARATOR + strconv.Itoa(user_id)).Err(); err != nil {
			log.Printf("delete password token for user [%d]failed: %s\n", user_id, err.Error())
		}
	}

	return &Token_info{
		Id:       int64(user_id),
		External: EXTERNAL_none,
	}, nil
}

func retrievePattern(pattern_key string, middle string) (user_id int, key string, value string, err error) {
	keys, err := DBs.Redis.Keys(pattern_key + REDIS_SEPARATOR + middle + REDIS_SEPARATOR + "*").Result()

	if err != nil {
		return 0, "", "", err
	}

	key_nbr := len(keys)

	if key_nbr != 1 {
		return 0, "", "", fmt.Errorf("expected one session for one token (%d finded)", key_nbr)
	}

	key = keys[0]

	user_id, err = strconv.Atoi(strings.SplitN(key, ".", 3)[2])

	if err != nil {
		return 0, key, "", fmt.Errorf("token backup is broken")
	}

	value, err = DBs.Redis.Get(keys[0]).Result()

	if err != nil {
		return 0, key, "", fmt.Errorf("token extraction failed")
	}

	return
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
