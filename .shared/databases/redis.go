package databases

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/.shared/environment"
)

const (
	REDIS_SEPARATOR                  = "."
	REDIS_EXTERNAL_none              = "none"
	REDIS_EXTERNAL_42                = "42"
	REDIS_EXTERNAL_google            = "google"
	REDIS_PATTERN_KEY_token          = "token"
	REDIS_PATTERN_KEY_search         = "search"
	REDIS_PATTERN_KEY_password_token = "password_token"
)

type RedisTokenInfo struct {
	Id       int64
	External string
}

var Redis *redis.Client

func AddToken(user_id int64, token string, external string) error {
	if len(external) == 0 {
		external = REDIS_EXTERNAL_none
	}

	err := Redis.Append(REDIS_PATTERN_KEY_token+REDIS_SEPARATOR+token+REDIS_SEPARATOR+strconv.Itoa(int(user_id)), external).Err()

	return err
}

func RetrieveToken(token string) (*RedisTokenInfo, error) {
	keys, err := Redis.Keys(REDIS_PATTERN_KEY_token + REDIS_SEPARATOR + token + REDIS_SEPARATOR + "*").Result()

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

	external, err := Redis.Get(keys[0]).Result()

	if err != nil {
		return nil, fmt.Errorf("token extraction failed")
	}

	return &RedisTokenInfo{
		Id:       int64(id),
		External: external,
	}, nil
}

func InitRedis() error {
	log.Println("start connection to redis on default address")
	Redis = redis.NewClient(&redis.Options{
		Addr:     environment.Redis.Host + ":" + environment.Redis.Port,
		Password: "",
		DB:       0,
	})

	if _, err := Redis.Ping().Result(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return nil
}
