package redisjson

import (
	"encoding/json"

	"github.com/pkg/errors"
	bridge "github.com/ridwanakf/url-shortener-service/internal"
)

type redisJson struct {
	bridge.Redis
	json bridge.Json
}

// NewRedisJson constructs a new RedisJson module which needs input of Redis and Json modules
func NewRedisJson(redis bridge.Redis, json bridge.Json) *redisJson {
	return &redisJson{Redis: redis, json: json}
}

func (rj *redisJson) GetUnmarshalled(key string, v interface{}) error {
	result, err := rj.Get(key)
	if err != nil {
		return errors.Wrap(err, "failed to get from redis")
	}

	err = rj.json.Unmarshal(result, v)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal redis result: %v", result)
	}

	return nil
}

func (rj *redisJson) SetexMarshalled(key string, seconds int, v interface{}) error {
	res, err := json.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal on key %s", key)
	}

	err = rj.Setex(key, seconds, string(res))
	if err != nil {
		return errors.Wrap(err, "failed to set redis")
	}

	return nil
}
