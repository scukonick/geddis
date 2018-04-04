package db

import (
	"testing"
	"time"
)

func TestStringStoreSetAndGet(t *testing.T) {
	value := "abcde"
	key := "zxc"

	store := newStringsStore(5)
	store.set(key, value, 0)
	res, err := store.get(key)
	if err != nil {
		t.Fatalf("want nil, got: %v", err)
	}

	if res != value {
		t.Errorf("Want: %s, got: %s", value, res)
	}
}

func TestStringStoreGetNotExistingElement(t *testing.T) {

	store := newStringsStore(5)
	_, err := store.get("zxc")
	if err != ErrNotFound {
		t.Fatalf("want Element not found, got: %v", err)
	}

}

func TestStringStoreSetAndDel(t *testing.T) {
	value := "abcde"
	key := "zxc"

	store := newStringsStore(5)
	store.set(key, value, 0)
	store.del(key)

	_, err := store.get(key)
	if err != ErrNotFound {
		t.Fatalf("want Element not found, got: %v", err)
	}
}

func TestStringStoreCleanExpired(t *testing.T) {
	value := "abcde"
	key := "zxc"

	store := newStringsStore(5)
	store.set(key, value, 1*time.Microsecond)

	time.Sleep(10 * time.Microsecond)

	_, err := store.get(key)
	if err != ErrNotFound {
		t.Fatalf("want Element not found, got: %v", err)
	}
}
