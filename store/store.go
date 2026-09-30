package store

// import "time"

//interface for defining a contract for the types of the stores like
//kv & TTL
type Storer interface {
	Set(key string, val string) error
	Get(key string) (string, error)
	Keys() []string
	Len() int
	// SetKeyWithEncryption(key string, val string) (string, error)
	Delete(key string) error
	Rename(oldKey, newKey string) error
	Clear() error
	Exists(key string) bool
	Pop(key string) (string, error)
	SetKeyWithEncryption(key, val string) (string, error)  

	//ttl methods 
	// Expire(key string,ttl time.Duration) error
	// TTL(key string) (time.Duration, error)
	// Persist(key string) error

	//atomic ops
	// Incr(key string) (int64, error)
	// Decr(key string) (int64, error)
}