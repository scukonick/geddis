package db

import (
	"testing"
	"time"
)

func TestStringStoreSetAndGet(t *testing.T) {
	value := "abcde"
	key := "zxc"

	store := NewGeddisStore(5)
	store.set(key, value, 0)
	res, err := store.GetStr(key)
	if err != nil {
		t.Fatalf("want nil, got: %v", err)
	}

	if res != value {
		t.Errorf("Want: %s, got: %s", value, res)
	}
}

func TestStringStoreSetAndGetSlice(t *testing.T) {
	value := []string{"abc", "xyz"}
	key := "zxc"

	store := NewGeddisStore(5)
	store.set(key, value, 0)
	res, err := store.GetArr(key)
	if err != nil {
		t.Fatalf("want nil, got: %v", err)
	}

	for i, val := range res {
		if value[i] != val {
			t.Errorf("Want %s, got: %s", val, value[i])
		}
	}

	// checking that value is returned by copy
	// not by pointer
	res[0] = "qwerty"
	if value[0] == res[0] {
		t.Errorf("expected abc, got %s", res[0])
	}
}

func TestStringStoreSetAndGetMap(t *testing.T) {
	value := map[string]string{
		"abc": "111",
		"xyz": "222",
	}

	key := "zxc"

	store := NewGeddisStore(5)
	store.set(key, value, 0)
	res, err := store.GetMap(key)
	if err != nil {
		t.Fatalf("want nil, got: %v", err)
	}

	for k := range res {
		if res[k] != value[k] {
			t.Errorf("Want %s, got: %s", value[k], res[k])
		}
	}

	// checking that value is returned by copy
	// not by pointer
	res["abc"] = "qwerty"
	if value["abc"] == res["abc"] {
		t.Errorf("expected 111, got %s", value["abc"])
	}
}

func TestStringStoreGetNotExistingElement(t *testing.T) {

	store := NewGeddisStore(5)
	_, err := store.GetStr("zxc")
	if err != ErrNotFound {
		t.Fatalf("want Element not found, got: %v", err)
	}

}

func TestStringStoreSetAndDel(t *testing.T) {
	value := "abcde"
	key := "zxc"

	store := NewGeddisStore(5)
	store.set(key, value, 0)
	store.Del(key)

	_, err := store.GetStr(key)
	if err != ErrNotFound {
		t.Fatalf("want Element not found, got: %v", err)
	}
}

func TestStringStoreCleanExpired(t *testing.T) {
	value := "abcde"
	key := "zxc"

	store := NewGeddisStore(5)
	store.set(key, value, 1*time.Microsecond)

	time.Sleep(10 * time.Microsecond)

	_, err := store.GetStr(key)
	if err != ErrNotFound {
		t.Fatalf("want Element not found, got: %v", err)
	}
}
func TestGetInvalidType(t *testing.T) {
	value := "abc"
	key := "xyz"

	store := NewGeddisStore(5)
	store.set(key, value, 0)

	_, err := store.GetMap(key)
	if err != ErrInvalidType {
		t.Errorf("Want Invalid Type Error, got: %v", err)
	}
}

func TestStorageKeys(t *testing.T) {
	keys := []string{"aaaa", "aaabbbb", "ccccc"}
	store := NewGeddisStore(5)

	for _, k := range keys {
		store.set(k, "123", 0)
	}

	storedKeys := store.Keys("aaa")
	if len(storedKeys) != 2 {
		t.Fatalf("expected 2, got: %d", len(storedKeys))
	}

	found := 0
	for _, k := range storedKeys {
		for _, j := range keys {
			if k == j {
				found++
			}
		}
	}
	if found != 2 {
		t.Errorf("Want 2, got: %d", found)
	}
}

func TestStorageGetByKey(t *testing.T) {
	val := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}

	key := "key"

	store := NewGeddisStore(5)
	store.SetMap(key, val, 0)

	res, err := store.GetByKey(key, "key1")
	if err != nil {
		t.Errorf("Want nil, got: %v", err)
	}

	if res != val["key1"] {
		t.Errorf("Want: %s, got: %s", val["key1"], res)
	}
}

func TestStorageGetByIndex(t *testing.T) {
	val := []string{"aaaa", "bbbb", "ccc"}

	key := "key"

	store := NewGeddisStore(5)
	store.SetArr(key, val, 0)

	res, err := store.GetByIndex(key, 1)
	if err != nil {
		t.Errorf("Want nil, got: %v", err)
	}

	if res != val[1] {
		t.Errorf("Want: %s, got: %s", val[1], res)
	}
}
