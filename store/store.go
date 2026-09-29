package store 


//interface for defining a contract for the types of the stores like 
//kv & TTL 
type Storer interface {
	Set(key string, val string) error
	Get(key string) (string, error)
	Keys() []string
	Len() int
	// SetKeyWithEncryption(key string, val string) (string, error)
	// Delete(key string)
	// Rename(oldKey, newKey string) error
	// Clear()
	// Count() int
	// Exists(key string) bool
	// Pop(key string) (string, bool)
}