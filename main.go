package main

import (
	"encoding/base64"
	"fmt"
	"simple-redis/store"
	"simple-redis/ttl"
	"time"

	"github.com/charmbracelet/log"
)

func SetKeyWithEncryption(s store.Storer, key, val string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(val))

	if err := s.Set(key, encoded); err != nil {
		return "", err
	}

	return encoded, nil
}

func main() {
	fmt.Println("Hello, World!")

	logger := log.NewWithOptions(nil, log.Options{
		ReportTimestamp: true,
	})

	// KV store
	kvStore := store.NewStore(4)

	// Wrap the KV store with logging.
	loggingStore := NewLoggingMiddleware(kvStore, logger)

	loggingStore.Set("a", "hello")

	value, err := loggingStore.Get("a")
	if err != nil {
		logger.Error("failed to get value", "error", err)
	} else {
		fmt.Println("Value:", value)
	}

	// TTL store
	ttlStore := ttl.NewTTLStore(5 * time.Second)

	err = ttlStore.SetWithTTL("e", "5678", 5*time.Second)
	if err != nil {
		logger.Error("SET failed", "key", "e", "error", err)
		return
	}
	ttlStore.SetWithTTL("f", "9012", 10*time.Second)

	time.Sleep(1 * time.Second)
	val, err := ttlStore.Get("e")
	if err != nil {
		logger.Error(
			"GET failed",
			"key", "e",
			"error", err,
		)
		return
	}

	fmt.Println("TTL value:", val)

	remaining, err := ttlStore.TTL("e")
	if err != nil {
		logger.Error(
			"TTL check failed",
			"key", "e",
			"error", err,
		)
		return
	}

	logger.Info(
		"TTL remaining",
		"key", "e",
		"remaining", remaining,
	)

	// Encryption
	encryptedVal, err := SetKeyWithEncryption(
		ttlStore,
		"c",
		"secretValue",
	)

	if err != nil {
		logger.Error(
			"failed to set encrypted value",
			"key",
			"c",
			"error",
			err,
		)
	} else {
		fmt.Println("Encrypted value:", encryptedVal)
	}

	// These calls bypass the logging middleware.
	kvStore.Set("b", "3444")
	kvStore.Set("d", "3444")

	val, err = kvStore.Get("a")
	if err != nil {
		logger.Error("failed to get KV value", "key", "a", "error", err)
	} else {
		fmt.Println("Value of a:", val)
	}

	if err := kvStore.Set("f", "4555"); err != nil {
		logger.Error("failed to set KV value", "key", "f", "error", err)
	}
}
