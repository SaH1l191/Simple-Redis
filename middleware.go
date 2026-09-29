package main

import (
	"simple-redis/store"

	"github.com/charmbracelet/log"
)

// LoggingMiddleware is a decorator that adds logging
// around operations performed on a Storer.
//
// LoggingMiddleware
//        |
//        | calls
//        ↓
//    store.Storer
//        |
//        ├── Store       → KV logic
//        |
//        └── TTLStore    → TTL logic
type LoggingMiddleware struct {
	store  store.Storer
	logger *log.Logger
}

func NewLoggingMiddleware(s store.Storer, logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		store:  s,
		logger: logger,
	}
}

func (m *LoggingMiddleware) Set(key string, val string) error {
	m.logger.Info("SET", "key", key)

	err := m.store.Set(key, val)
	if err != nil {
		m.logger.Error("SET failed", "key", key, "error", err)
		return err
	}

	m.logger.Info("SET successful", "key", key)

	return nil
}

func (m *LoggingMiddleware) Get(key string) (string, error) {
	m.logger.Info("GET", "key", key)

	val, err := m.store.Get(key)
	if err != nil {
		m.logger.Error("GET failed", "key", key, "error", err)
		return "", err
	}

	m.logger.Info("GET successful", "key", key)

	return val, nil
}

func (m *LoggingMiddleware) Keys() []string {
	m.logger.Info("KEYS")

	keys := m.store.Keys()

	m.logger.Info("KEYS successful", "count", len(keys))

	return keys
}

func (m *LoggingMiddleware) Len() int {
	m.logger.Info("LEN")

	length := m.store.Len()

	m.logger.Info("LEN successful", "count", length)

	return length
}
