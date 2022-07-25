package queries

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/trixky/hypertube/.shared/databases"
)

func CacheSearch(path *string, results *string) error {
	// Save in redis and set the ttl to 5min
	if err := databases.Redis.Set(databases.REDIS_PATTERN_KEY_search+*path, *results, 5*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}

func RetrieveSearch(path *string) (string, error) {
	// Check if the search exists on redis
	results, err := databases.Redis.Get(databases.REDIS_PATTERN_KEY_search + *path).Result()
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
