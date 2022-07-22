package databases

import (
	"errors"
	"strconv"

	"github.com/trixky/hypertube/api-auth/databases"
)

type RedisDatabaseMock struct {
	Objects map[string]interface{}
}

var InitialsRedisMockObjects = map[string]interface{}{
	databases.REDIS_PATTERN_KEY_password_token + databases.REDIS_SEPARATOR + "c57a31ea-0125-4a9e-b602-fb9faf7e3f45" + databases.REDIS_SEPARATOR + "0": databases.EXTERNAL_none,
	databases.REDIS_PATTERN_KEY_password_token + databases.REDIS_SEPARATOR + "eeee318f-1b68-49f8-9ad7-e8695ad114a9" + databases.REDIS_SEPARATOR + "1": databases.EXTERNAL_none,
}

// ---------------------------------- INIT

// InitRedisMock intisalizes the redis (mock)
func InitRedisMock() {
	databases.DBs.RedisDatabase = nil
	databases.DBs.RedisQueries = RedisDatabaseMock{
		Objects: InitialsRedisMockObjects,
	}
}

// ---------------------------------- IMPLEMENTATION

// AddToken (mock)
func (rdm RedisDatabaseMock) AddToken(user_id int64, token string, external string) error {
	if len(external) == 0 {
		// If no external info is specified
		external = databases.EXTERNAL_none
	}

	rdm.Objects[databases.REDIS_PATTERN_KEY_token+databases.REDIS_SEPARATOR+token+databases.REDIS_SEPARATOR+strconv.Itoa(int(user_id))] = external

	return nil
}

// RetrieveToken (mock)
func (rdm RedisDatabaseMock) RetrieveToken(token string) (*databases.RedisTokenInfo, error) {
	key_start := databases.REDIS_PATTERN_KEY_token + databases.REDIS_SEPARATOR + token + databases.REDIS_SEPARATOR

	for i := 0; i < 30; i++ {
		if value, ok := rdm.Objects[key_start+strconv.Itoa(i)]; ok {
			return &databases.RedisTokenInfo{
				Id:       int64(i),
				External: value.(string),
			}, nil
		}
	}

	return nil, errors.New("no token retrieve")
}

// AddPasswordToken (mock)
func (rdm RedisDatabaseMock) AddPasswordToken(user_id int64, token string) error {
	rdm.Objects[databases.REDIS_PATTERN_KEY_password_token+databases.REDIS_SEPARATOR+token+databases.REDIS_SEPARATOR+strconv.Itoa(int(user_id))] = ""

	return nil
}

// RetrievePasswordToken (mock)
func (rdm RedisDatabaseMock) RetrievePasswordToken(token string, delete_after bool) (*databases.RedisTokenInfo, error) {
	key_start := databases.REDIS_PATTERN_KEY_password_token + databases.REDIS_SEPARATOR + token + databases.REDIS_SEPARATOR

	for i := 0; i < 30; i++ {
		if value, ok := rdm.Objects[key_start+strconv.Itoa(i)]; ok {
			if delete_after {
				delete(rdm.Objects, "sdf")
			}

			return &databases.RedisTokenInfo{
				Id:       int64(i),
				External: value.(string),
			}, nil
		}
	}

	return nil, errors.New("no token retrieve")
}
