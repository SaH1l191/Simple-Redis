package store

import (
	"errors"
	"slices"
	"testing"
)

func TestKeys_ReturnAllKeysSorted(t *testing.T) {
	store := NewStore()
	store.Set("a", "1")
	store.Set("b", "2")
	store.Set("c", "234")

	got := store.Keys()
	want := []string{"a", "b", "c"}
	if !slices.Equal(want, got) {
		t.Errorf("Key () = %v,want %v", got, want)
	}
}

func Test_EmptyStore(t *testing.T) {
	store := NewStore()

	got := store.Keys()
	if len(got) != 0 {
		t.Errorf("Keys() on Empty store = %v, wanted empty slice", got)
	}
}

func Test_SetGe_EmptyKeys(t *testing.T) {
	store := NewStore()

	if _, err := store.Get(""); err == nil || !errors.Is(err, ErrEmptyKey) {
		t.Error("Get() Failed")
	}
}

func Test_SetGet(t *testing.T) {
	store := NewStore()
	store.Set("a", "1")

	val, err := store.Get("a")
	if err != nil {
		t.Errorf("Get() Failed: %v", err)
	}
	if val != "1" {
		t.Errorf("Get() Failed: expected '1', got '%s'", val)
	}
	if val, err := store.Get("missing"); err == nil || val != "" {
		t.Error("Get() Failed")
	}
}
