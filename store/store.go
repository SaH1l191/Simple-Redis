package store

import (
	"errors"
	"fmt"
	"sort"
)

type Store struct {
	data map[string]string
}

var ErrEmptyKey = errors.New("Key cannot be empty")

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *Store) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	val, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

func (s *Store) Set(key string, val string) {
	s.data[key] = val
}

func (s *Store) Delete(key string) {
	delete(s.data, key)
}

func (s *Store) Rename(oldKey string, newKey string) error {
	val, err := s.Get(oldKey)
	if err != nil {
		return err
	}
	delete(s.data, oldKey)
	s.Set(newKey, val)
	return nil
}

func (s *Store) Clear() {
	s.data = make(map[string]string)
}

func (s *Store) Count() int {
	return len(s.data)
}

func (s *Store) Exists(key string) bool {
	_, exists := s.data[key]
	return exists
}

func (s *Store) Pop(key string) (string, bool) {
	val, ok := s.Get(key)
	if ok != nil {
		return "", false
	}
	delete(s.data, key)
	return val, true
}
