package store

import (
	"encoding/base64"
	"errors"
	// "fmt"
	"sort"
	"sync"
)

//instead of having a TTL for each key , we have a explicit TTL expiry time for each key -> memory Efficient and fast lookup for TTLs.
//atomic expire+delete by janitor delete from both atomically
// no value copy
//so fast TTL lookup ,
type Store struct {
	data    map[string]string
	maxSize int // 0 means no limit
	mu sync.RWMutex 
	// expiry   map[string]time.Time
}

var (
	ErrEmptyKey  = errors.New("key cannot be empty")
	ErrStoreFull = errors.New("store is full")
	ErrKeyNotFound = errors.New("key not found")
)

func NewStore(maxSize int) *Store {
	return &Store{
		data:    make(map[string]string),
		maxSize: maxSize,
	}
}

func (s *Store) Set(key string, val string) error {
	if key == "" {
		return ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]

	if s.maxSize > 0 && len(s.data) >= s.maxSize && !exists {
		return ErrStoreFull
	}

	s.data[key] = val

	return nil
}

func (s *Store) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	s.mu.RLock()

	val, exists := s.data[key]

	s.mu.RUnlock()
	if !exists {
		return "", ErrKeyNotFound
	}
	

	return val, nil
}

func (s *Store) Keys() []string {

	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Store) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	

	s.mu.Lock()
	defer s.mu.Unlock()

	if _,ok := s.data[key]; !ok {
		return ErrKeyNotFound
	}

	delete(s.data, key)
	return nil

}

func (s *Store) Rename(oldKey string, newKey string) error {

	if oldKey == "" || newKey == "" {
		return ErrEmptyKey
	}
	if oldKey == newKey {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	val, err := s.Get(oldKey)
	if err != nil {
		return err
	}

	if err := s.Set(newKey, val); err != nil {
		return err
	}

	delete(s.data, oldKey)

	return nil
}

func (s *Store) Clear() error{
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]string)
	return nil
}

func (s *Store) Count() int {
	return len(s.data)
}

func (s *Store) Exists(key string) bool {
	_, exists := s.Get(key)

	return exists == nil
}

func (s *Store) Pop(key string) (string, error) {

	if key == "" {
		return "", ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	val, err := s.Get(key)
	if err != nil {
		return "", ErrKeyNotFound
	}
	delete(s.data, key)

	return val, nil
}


func (s *Store) SetKeyWithEncryption(key, val string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(val))
	if err := s.Set(key, encoded); err != nil {
		return "", err
	}
	return encoded, nil
}

// func (s *Store) Expire(key string, ttl time.Duration) error {
// 	if key == "" {
// 		return ErrEmptyKey
// 	}

// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if _, ok := s.data[key]; !ok {
// 		return ErrKeyNotFound
// 	}

// 	if ttl <= 0 {
// 		delete(s.expiry, key)
// 		return nil
// 	}

// 	s.expiry[key] = time.Now().Add(ttl)
// 	return nil
// }


// func (s *Store) TTL(key string) (time.Duration, error) {
// 	if key == "" {
// 		return 0, ErrEmptyKey
// 	}

// 	s.mu.RLock()
// 	exp, hasExp := s.expiry[key]
// 	_, ok := s.data[key]
// 	s.mu.RUnlock()

// 	if !ok {
// 		return 0, ErrKeyNotFound
// 	}

// 	if !hasExp {
// 		return -1, nil // No expiry set
// 	}

// 	remaining := time.Until(exp)
// 	if remaining <= 0 {
// 		return -2, nil // Expired
// 	}

// 	return remaining, nil
// }

// // Persist removes expiration from a key
// func (s *Store) Persist(key string) error {
// 	if key == "" {
// 		return ErrEmptyKey
// 	}

// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if _, ok := s.data[key]; !ok {
// 		return ErrKeyNotFound
// 	}

// 	delete(s.expiry, key)
// 	return nil
// }

// // ==================== Atomic Counters ====================

// // Incr increments the integer value at key by 1
// func (s *Store) Incr(key string) (int64, error) {
// 	return s.IncrBy(key, 1)
// }

// // Decr decrements the integer value at key by 1
// func (s *Store) Decr(key string) (int64, error) {
// 	return s.IncrBy(key, -1)
// }

// // IncrBy increments the integer value at key by delta
// func (s *Store) IncrBy(key string, delta int64) (int64, error) {
// 	if key == "" {
// 		return 0, ErrEmptyKey
// 	}

// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	valStr, ok := s.data[key]
// 	var current int64
// 	if ok {
// 		var parseErr error
// 		current, parseErr = parseInt(valStr)
// 		if parseErr != nil {
// 			return 0, fmt.Errorf("value is not an integer: %s", valStr)
// 		}
// 	}

// 	newVal := current + delta
// 	s.data[key] = fmt.Sprintf("%d", newVal)
// 	return newVal, nil
// }

// // parseInt parses a string to int64
// func parseInt(s string) (int64, error) {
// 	var n int64
// 	var neg bool
// 	for i, c := range s {
// 		if i == 0 && c == '-' {
// 			neg = true
// 			continue
// 		}
// 		if c < '0' || c > '9' {
// 			return 0, fmt.Errorf("invalid integer: %s", s)
// 		}
// 		n = n*10 + int64(c-'0')
// 	}
// 	if neg {
// 		n = -n
// 	}
// 	return n, nil
// }

// // ==================== Internal: Expiry Janitor ====================

// // startExpiryJanitor runs a background goroutine to clean expired keys
// func (s *Store) startExpiryJanitor() {
// 	ticker := time.NewTicker(100 * time.Millisecond)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		s.cleanExpired()
// 	}
// }

// // cleanExpired removes expired keys (called periodically)
// func (s *Store) cleanExpired() {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	now := time.Now()
// 	for k, exp := range s.expiry {
// 		if now.After(exp) {
// 			delete(s.data, k)
// 			delete(s.expiry, k)
// 		}
// 	}
// }

// func (s *Store) SetKeyWithEncryption(key, val string) (string, error) {
// 	encoded := base64.StdEncoding.EncodeToString([]byte(val))
// 	if err := s.Set(key, encoded); err != nil {
// 		return "", err
// 	}
// 	return encoded, nil
// }

// // Add stub TTL methods (return "not implemented" or delegate)
// func (s *Store) SetWithTTL(key, val string, ttl time.Duration) error {
// 	return s.Set(key, val) // Ignore TTL for basic store
// }

// func (s *Store) TTL(key string) (time.Duration, error) {
// 	return 0, ErrNotSupported
// }