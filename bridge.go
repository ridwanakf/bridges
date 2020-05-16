package go_bridges

// Json module used for encoding and decoding JSON string to Golang's struct
//go:generate mockgen -destination=json/jsoniterator_mock.go -package=json github.com/ridwanakf/go-bridges Json
type Json interface {
	// Unmarshal decodes JSON string to a struct
	// v should be a pointer
	Unmarshal(data []byte, v interface{}) error

	// Marshal encodes struct into JSON
	// v should not be a pointer
	Marshal(v interface{}) ([]byte, error)
}

// Redis module is a client to connect to redis
//go:generate mockgen -destination=redis/redigo_mock.go -package=redis github.com/ridwanakf/go-bridges Redis
type Redis interface {
	// Get value of the specified key. Will return ErrNil if the return value is nil
	Get(key string) ([]byte, error)
	// Set sets value to key in redis without any additional options
	Set(key, value string) error
	Setex(key string, seconds int, value string) error
	// Setnx sets a value to a key with specified timeouts. Will return false if the key exists
	Setnx(key string, seconds int, value string) (bool, error)
	// HMGet gets a value of multiple fields from hash key. Will return ErrNil if the return value is nil
	HMGet(key string, fields ...string) ([][]byte, error)
	Exists(key string) (bool, error)
	Expire(key string, seconds int) (bool, error)
	ExpireAt(key string, timestamp int64) (bool, error)
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	TTL(key string) (int64, error)
	HGet(key string, field string) ([]byte, error)
	HExists(key string, field string) (bool, error)
	HGetAll(key string) (map[string]string, error)
	HSet(key string, field string, value string) (bool, error)
	HKeys(key string) ([]string, error)
	HDel(key string, fields ...string) (int64, error)
	Del(key ...interface{}) (int64, error)

	// Close releases all the connections and resources to redis
	Close() error
}

// RedisJson adds the capability to Redis by adding capabilities to Unmarshal and Marshal JSON on this very module
//go:generate mockgen -destination=redisjson/redisjson_mock.go -package=redisjson github.com/ridwanakf/go-bridges RedisJson
type RedisJson interface {
	Redis
	// GetUnmarshalled fetches value from redis key
	// v should be pointer
	GetUnmarshalled(key string, v interface{}) error

	// SetexMarshalled sets value with expiry
	// v should not be a pointer
	SetexMarshalled(key string, seconds int, v interface{}) error
}
