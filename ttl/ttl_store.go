package ttl

import (
	"fmt"
	"simple-redis/store"
	"sort"
	"time"
)

type ttlEntry struct {
	value     string
	expiresAt time.Time
}

type TTLStore struct {
	data       map[string]ttlEntry
	defaultTTL time.Duration
}

func NewTTLStore(defaultTTL time.Duration) *TTLStore {
	return &TTLStore{
		data:       make(map[string]ttlEntry),
		defaultTTL: defaultTTL,
	}
}

func (s *TTLStore) isExpired(entry ttlEntry) bool {
	return !time.Now().Before(entry.expiresAt)
}

func (s *TTLStore) Set(key, value string) error {
	return s.SetWithTTL(key, value, s.defaultTTL)
}

func (s *TTLStore) SetWithTTL(
	key string, value string, ttl time.Duration,
) error {
	if key == "" {
		return store.ErrEmptyKey
	}

	if ttl <= 0 {
		return fmt.Errorf("TTL must be greater than zero")
	}

	s.data[key] = ttlEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (s *TTLStore) Get(key string) (string, error) {
	if key == "" {
		return "", store.ErrEmptyKey
	}
	entry, exists := s.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}

	if s.isExpired(entry) {
		delete(s.data, key)
		return "", fmt.Errorf("key %s has expired", key)
	}
	return entry.value, nil
}

func (s *TTLStore) Keys() []string {
	keys := make([]string, 0, len(s.data))

	for key, entry := range s.data {
		if s.isExpired(entry) {
			delete(s.data, key)
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (s *TTLStore) Delete(key string) {
	delete(s.data, key)
}

func (s *TTLStore) Rename(oldKey, newKey string) error {
	entry, exists := s.data[oldKey]
	if !exists {
		return fmt.Errorf("key not found: %s", oldKey)
	}

	if s.isExpired(entry) {
		delete(s.data, oldKey)
		return fmt.Errorf("key %s has expired", oldKey)
	}

	if oldKey == newKey {
		return nil
	}

	s.data[newKey] = entry
	delete(s.data, oldKey)

	return nil
}

func (s *TTLStore) Clear() {
	s.data = make(map[string]ttlEntry)
}

func (s *TTLStore) Count() int {
	count := 0

	for key, entry := range s.data {
		if s.isExpired(entry) {
			delete(s.data, key)
			continue
		}

		count++
	}
	return count
}

func (s *TTLStore) Exists(key string) bool {
	entry, exists := s.data[key]
	if !exists {
		return false
	}

	if s.isExpired(entry) {
		delete(s.data, key)
		return false
	}

	return true
}

func (s *TTLStore) Pop(key string) (string, bool) {
	entry, exists := s.data[key]
	if !exists {
		return "", false
	}

	if s.isExpired(entry) {
		delete(s.data, key)
		return "", false
	}

	delete(s.data, key)
	return entry.value, true
}

func (s *TTLStore) TTL(key string) (time.Duration, error) {
	entry, exists := s.data[key]
	if !exists {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	if s.isExpired(entry) {
		delete(s.data, key)
		return 0, fmt.Errorf("key %s has expired", key)
	}

	return time.Until(entry.expiresAt), nil
}

func (s *TTLStore) Len() int {
	count := 0

	for key, entry := range s.data {
		if s.isExpired(entry) {
			delete(s.data, key)
			continue
		}
		count++
	}
	return count
}
