package store

import (
	"errors"
	"fmt"
	"sort"
)

type Store struct {
	data    map[string]string
	maxSize int // 0 means no limit
}

var (
	ErrEmptyKey  = errors.New("key cannot be empty")
	ErrStoreFull = errors.New("store is full")
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

	_, exists := s.data[key]

	if s.maxSize > 0 && s.Len() >= s.maxSize && !exists {
		return ErrStoreFull
	}

	s.data[key] = val

	return nil
}

func (s *Store) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	val, exists := s.data[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return val, nil
}

func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (s *Store) Len() int {
	return len(s.data)
}

func (s *Store) Delete(key string) {
	delete(s.data, key)
}

func (s *Store) Rename(oldKey string, newKey string) error {
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
	val, err := s.Get(key)
	if err != nil {
		return "", false
	}

	delete(s.data, key)

	return val, true
}
