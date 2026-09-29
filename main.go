package main

import (
	"encoding/base64"
	"fmt"
	"simple-redis/store"
	"simple-redis/ttl"
	"time"
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
	kvStore := store.NewStore(4)
	ttlStore := ttl.NewTTLStore(time.Second * 5)

	kvStore.Set("a", "23")
	kvStore.Set("b", "3444")

	ttlStore.SetWithTTL("e", "5678", time.Second*5)
	ttlStore.SetWithTTL("f", "9012", time.Second*10)

	time.Sleep(1 * time.Second)

	if val, err := ttlStore.Get("e"); err != nil {
		fmt.Println("Error getting key 'e':", err)
	} else {
		fmt.Printf("%v\n",val)
	}

	encryptedVal, err := SetKeyWithEncryption(ttlStore, "c", "secretValue")
	if err != nil {
		fmt.Println("Error setting key with encryption:", err)
	} else {
		fmt.Println("Encrypted value for key 'c':", encryptedVal)
	}
	kvStore.Set("d", "3444")

	val, err := kvStore.Get("a")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Value of a is ", val)

	if err := kvStore.Set("f", "4555"); err != nil {
		fmt.Printf("Error %v", err)
	}

}
