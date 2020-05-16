package redis

import (
	"log"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// redigo holds the pool of redigo for redis operations
type redigo struct {
	pool *redis.Pool
}

type ConfigOptions struct {
	Address   string
	Timeout   int
	MaxIdle   int
	MaxActive int
}

// NewRedigo constructs a new Redis-client using Redigo library
func NewRedigo(redisConfig ConfigOptions) *redigo {
	// Reference: https://github.com/pete911/examples-redigo
	pool := &redis.Pool{

		MaxIdle:     redisConfig.MaxIdle,
		MaxActive:   redisConfig.MaxActive,
		IdleTimeout: time.Duration(redisConfig.Timeout) * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisConfig.Address)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	mustTryDialling(pool)

	return &redigo{pool: pool}
}

func mustTryDialling(pool *redis.Pool) {
	if _, err := pool.Dial(); err != nil {
		if closeErr := pool.Close(); closeErr != nil {
			log.Println("Failed to close Redis connection pool", closeErr)
		}
		log.Fatalln("Failed to dial Redis connection", err)
	}
}

// Get gets the value from redis in []byte form
func (r *redigo) Get(key string) ([]byte, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Bytes(con.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return nil, err
	} else if err == redis.ErrNil {
		return nil, ErrNil
	}

	return data, nil
}

// Set sets value to key in redis without any additional options
func (r *redigo) Set(key, value string) error {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.String(con.Do("SET", key, value))
	if err != nil && err != redis.ErrNil {
		return err
	}
	//extra check for set operations
	if err == redis.ErrNil || !strings.EqualFold("OK", data) {
		return ErrFailed
	}

	return nil
}

// Setex sets the value to a key with timeout in seconds
func (r *redigo) Setex(key string, seconds int, value string) error {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.String(con.Do("SET", key, value, "ex", seconds))
	if err != nil && err != redis.ErrNil {
		return err
	}
	//extra check for set operations
	if err == redis.ErrNil || !strings.EqualFold("OK", data) {
		return ErrFailed
	}

	return nil
}

// Setnx sets a value to a key with specified timeouts. Will return false if the key exists
func (r *redigo) Setnx(key string, seconds int, value string) (bool, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.String(con.Do("SET", key, value, "ex", seconds, "nx"))
	if err != nil && err != redis.ErrNil {
		return false, err
	}
	//extra check for set operations
	if err == redis.ErrNil || !strings.EqualFold("OK", data) {
		return false, nil
	}

	return true, nil
}

// HMGet gets a value of multiple fields from hash key
func (r *redigo) HMGet(key string, fields ...string) ([][]byte, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.ByteSlices(con.Do("HMGET", key, fields))
	if err != nil && err != redis.ErrNil {
		return nil, err
	} else if err == redis.ErrNil {
		return nil, ErrNil
	}

	return data, nil
}

// Exists checks whether the key exists in redis
func (r *redigo) Exists(key string) (bool, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int(con.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	if data != 1 {
		return false, nil
	}

	return true, nil
}

// Expire sets the ttl of a key to specified value in seconds
func (r *redigo) Expire(key string, seconds int) (bool, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int(con.Do("EXPIRE", key, seconds))
	if err != nil {
		return false, err
	}

	if data != 1 {
		return false, nil
	}

	return true, nil
}

// ExpireAt sets the ttl of a key to a certain timestamp
func (r *redigo) ExpireAt(key string, timestamp int64) (bool, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int(con.Do("EXPIREAT", key, timestamp))
	if err != nil {
		return false, err
	}

	if data != 1 {
		return false, nil
	}

	return true, nil
}

// Incr increments the integer value of a key by 1
func (r *redigo) Incr(key string) (int64, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int64(con.Do("INCR", key))
	if err != nil {
		return 0, err
	}

	return data, nil
}

// Decr decrements the integer value of a key by 1
func (r *redigo) Decr(key string) (int64, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int64(con.Do("DECR", key))
	if err != nil {
		return 0, err
	}

	return data, nil
}

// TTL gets the time to live of a key / expiry time
func (r *redigo) TTL(key string) (int64, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int64(con.Do("TTL", key))
	if err != nil {
		return 0, err
	}

	return data, nil
}

// HGet gets the value of a hash field
func (r *redigo) HGet(key string, field string) ([]byte, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Bytes(con.Do("HGET", key, field))
	if err != nil && err != redis.ErrNil {
		return data, err
	}

	return data, nil
}

// HExists determines if a hash field exists
func (r *redigo) HExists(key string, field string) (bool, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int(con.Do("HEXISTS", key, field))
	if err != nil {
		return false, err
	}

	if data != 1 {
		return false, nil
	}

	return true, nil
}

// HGetAll gets all the fields and values in a hash
func (r *redigo) HGetAll(key string) (map[string]string, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.StringMap(con.Do("HGETALL", key))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	return data, nil
}

// HSet sets the string value of a hash field
func (r *redigo) HSet(key string, field string, value string) (bool, error) {
	con := r.pool.Get()
	defer con.Close()

	res, err := redis.Int(con.Do("HSET", key, field, value))
	if err != nil {
		return false, err
	}

	// 1 if field is a new field in the hash and value was set.
	// 0 if field already exists in the hash and the value was updated.
	if res != 1 && res != 0 {
		return false, nil
	}

	return true, nil
}

// HKeys gets all the fields in a hash
func (r *redigo) HKeys(key string) ([]string, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Strings(con.Do("HKEYS", key))
	if err != nil {
		return nil, err
	}

	return data, nil
}

// HDel deletes a hash field
func (r *redigo) HDel(key string, fields ...string) (int64, error) {
	con := r.pool.Get()
	defer con.Close()

	params := make([]interface{}, 0, len(fields)+1)

	params = append(params, key)
	for _, f := range fields {
		params = append(params, f)
	}

	data, err := redis.Int64(con.Do("HDEL", params...))
	if err != nil {
		return 0, ErrFailed
	}

	return data, nil
}

// Del deletes a key
func (r *redigo) Del(key ...interface{}) (int64, error) {
	con := r.pool.Get()
	defer con.Close()

	data, err := redis.Int64(con.Do("DEL", key...))
	if err == nil && data > 0 {
		keyString := make([]string, len(key))
		for i, v := range key {
			keyString[i] = v.(string)
		}
	}

	return data, err
}

func (r *redigo) Close() error {
	return r.pool.Close()
}
