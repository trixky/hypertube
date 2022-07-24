package databases

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/.shared/environment"
)

type RedisDatabase struct{}

const (
	REDIS_EXTERNAL_google            = "google"
	REDIS_SEPARATOR                  = "."
	REDIS_PATTERN_KEY_token          = "token"
	REDIS_PATTERN_KEY_password_token = "password_token"
)

// ErrorIsDuplication checks if an sql error is an duplication error
func ErrorIsDuplication(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}

// ---------------------------------- INIT

// InitRedis intisalizes the redis connection
func InitRedis() error {
	DBs.RedisDatabase = redis.NewClient(&redis.Options{
		Addr:     environment.Redis.RedisHost + ":6379",
		Password: "",
		DB:       0,
	})

	// Test the connection with a ping
	if _, err := DBs.RedisDatabase.Ping().Result(); err != nil {
		return err
	}

	DBs.RedisQueries = RedisDatabase{}

	return nil
}

// ---------------------------------- IMPLEMENTATION

// AddToken adds/save an authentification token to redis
func (rd RedisDatabase) AddToken(user_id int64, token string, external string) error {
	if len(external) == 0 {
		// If no external info is specified
		external = EXTERNAL_none
	}

	err := DBs.RedisDatabase.Append(REDIS_PATTERN_KEY_token+REDIS_SEPARATOR+token+REDIS_SEPARATOR+strconv.Itoa(int(user_id)), external).Err()

	return err
}

// RetrieveToken retrieves/get an authentification token to redis
func (rd RedisDatabase) RetrieveToken(token string) (*RedisTokenInfo, error) {
	keys, err := DBs.RedisDatabase.Keys(REDIS_PATTERN_KEY_token + REDIS_SEPARATOR + token + REDIS_SEPARATOR + "*").Result()

	if err != nil {
		return nil, err
	}

	key_nbr := len(keys)

	if key_nbr != 1 {
		// If zero or too many token have been found
		return nil, fmt.Errorf("expected one session for one token (%d finded)", key_nbr)
	}

	// Extract the user id from the token
	user_id, err := strconv.Atoi(strings.SplitN(keys[0], ".", 3)[2])

	if err != nil {
		return nil, fmt.Errorf("token backup is broken")
	}

	// Extract the external info from the token
	external, err := DBs.RedisDatabase.Get(keys[0]).Result()

	if err != nil {
		return nil, fmt.Errorf("token extraction failed")
	}

	return &RedisTokenInfo{
		Id:       int64(user_id),
		External: external,
	}, nil
}

// AddPasswordToken adds/save an passowrd token to redis
func (rd RedisDatabase) AddPasswordToken(user_id int64, token string) error {
	err := DBs.RedisDatabase.Set(REDIS_PATTERN_KEY_password_token+REDIS_SEPARATOR+token+REDIS_SEPARATOR+strconv.Itoa(int(user_id)), "", 10*time.Minute).Err()

	return err
}

// RetrievePasswordToken retrieves/get an passowrd token to redis
func (rd RedisDatabase) RetrievePasswordToken(token string, delete_after bool) (*RedisTokenInfo, error) {
	// Retrieves the user id of the token if exists
	user_id, _, _, err := rd.retrievePattern(REDIS_PATTERN_KEY_password_token, token)

	if err != nil {
		return nil, err
	}

	// Delete the password token because is for single use
	if delete_after {
		if err := DBs.RedisDatabase.Del(REDIS_PATTERN_KEY_password_token + REDIS_SEPARATOR + token + REDIS_SEPARATOR + strconv.Itoa(user_id)).Err(); err != nil {
			log.Printf("delete password token for user [%d]failed: %s\n", user_id, err.Error())
		}
	}

	return &RedisTokenInfo{
		Id:       int64(user_id),
		External: EXTERNAL_none,
	}, nil
}

// retrievePattern retrieves/get any type of infos from a specific pattern
// Pattern: <pattern_key>.<middle>.<* (any user_id)>
func (rd RedisDatabase) retrievePattern(pattern_key string, middle string) (user_id int, key string, value string, err error) {
	keys, err := DBs.RedisDatabase.Keys(pattern_key + REDIS_SEPARATOR + middle + REDIS_SEPARATOR + "*").Result()

	if err != nil {
		return 0, "", "", err
	}

	key_nbr := len(keys)

	if key_nbr != 1 {
		// If zero or too many token have been found
		return 0, "", "", fmt.Errorf("expected one session for one token (%d finded)", key_nbr)
	}

	key = keys[0]

	// Extract the user id from the token
	user_id, err = strconv.Atoi(strings.SplitN(key, ".", 3)[2])

	if err != nil {
		return 0, key, "", fmt.Errorf("token backup is broken")
	}

	// Extract the value from the token
	value, err = DBs.RedisDatabase.Get(keys[0]).Result()

	if err != nil {
		return 0, key, "", fmt.Errorf("token extraction failed")
	}

	return
}
