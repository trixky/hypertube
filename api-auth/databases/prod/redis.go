package databases

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/api-auth/databases"
)

type RedisDatabase struct{}

// ---------------------------------- INIT

// InitRedis intisalizes the redis connection
func InitRedis() error {
	databases.DBs.RedisDatabase = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// Test the connection with a ping
	if _, err := databases.DBs.RedisDatabase.Ping().Result(); err != nil {
		return err
	}

	databases.DBs.RedisQueries = RedisDatabase{}

	return nil
}

// ---------------------------------- IMPLEMENTATION

// AddToken adds/save an authentification token to redis
func (rd RedisDatabase) AddToken(user_id int64, token string, external string) error {
	if len(external) == 0 {
		// If no external info is specified
		external = databases.EXTERNAL_none
	}

	err := databases.DBs.RedisDatabase.Append(databases.REDIS_PATTERN_KEY_token+databases.REDIS_SEPARATOR+token+databases.REDIS_SEPARATOR+strconv.Itoa(int(user_id)), external).Err()

	return err
}

// RetrieveToken retrieves/get an authentification token to redis
func (rd RedisDatabase) RetrieveToken(token string) (*databases.RedisTokenInfo, error) {
	keys, err := databases.DBs.RedisDatabase.Keys(databases.REDIS_PATTERN_KEY_token + databases.REDIS_SEPARATOR + token + databases.REDIS_SEPARATOR + "*").Result()

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
	external, err := databases.DBs.RedisDatabase.Get(keys[0]).Result()

	if err != nil {
		return nil, fmt.Errorf("token extraction failed")
	}

	return &databases.RedisTokenInfo{
		Id:       int64(user_id),
		External: external,
	}, nil
}

// AddPasswordToken adds/save an passowrd token to redis
func (rd RedisDatabase) AddPasswordToken(user_id int64, token string) error {
	err := databases.DBs.RedisDatabase.Set(databases.REDIS_PATTERN_KEY_password_token+databases.REDIS_SEPARATOR+token+databases.REDIS_SEPARATOR+strconv.Itoa(int(user_id)), "", 10*time.Minute).Err()

	return err
}

// RetrievePasswordToken retrieves/get an passowrd token to redis
func (rd RedisDatabase) RetrievePasswordToken(token string, delete_after bool) (*databases.RedisTokenInfo, error) {
	// Retrieves the user id of the token if exists
	user_id, _, _, err := rd.retrievePattern(databases.REDIS_PATTERN_KEY_password_token, token)

	if err != nil {
		return nil, err
	}

	// Delete the password token because is for single use
	if delete_after {
		if err := databases.DBs.RedisDatabase.Del(databases.REDIS_PATTERN_KEY_password_token + databases.REDIS_SEPARATOR + token + databases.REDIS_SEPARATOR + strconv.Itoa(user_id)).Err(); err != nil {
			log.Printf("delete password token for user [%d]failed: %s\n", user_id, err.Error())
		}
	}

	return &databases.RedisTokenInfo{
		Id:       int64(user_id),
		External: databases.EXTERNAL_none,
	}, nil
}

// retrievePattern retrieves/get any type of infos from a specific pattern
// Pattern: <pattern_key>.<middle>.<* (any user_id)>
func (rd RedisDatabase) retrievePattern(pattern_key string, middle string) (user_id int, key string, value string, err error) {
	keys, err := databases.DBs.RedisDatabase.Keys(pattern_key + databases.REDIS_SEPARATOR + middle + databases.REDIS_SEPARATOR + "*").Result()

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
	value, err = databases.DBs.RedisDatabase.Get(keys[0]).Result()

	if err != nil {
		return 0, key, "", fmt.Errorf("token extraction failed")
	}

	return
}
