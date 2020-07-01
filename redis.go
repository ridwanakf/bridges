package bridges

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
