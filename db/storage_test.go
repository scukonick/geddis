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
	res, err := store.getStr(key)
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

	store := newStringsStore(5)
	store.set(key, value, 0)
	res, err := store.getArr(key)
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

	store := newStringsStore(5)
	store.set(key, value, 0)
	res, err := store.getMap(key)
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

	store := newStringsStore(5)
	_, err := store.getStr("zxc")
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

	_, err := store.getStr(key)
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

	_, err := store.getStr(key)
	if err != ErrNotFound {
		t.Fatalf("want Element not found, got: %v", err)
	}
}
func TestGetInvalidType(t *testing.T) {
	value := "abc"
	key := "xyz"

	store := newStringsStore(5)
	store.set(key, value, 0)

	_, err := store.getMap(key)
	if err != ErrInvalidType {
		t.Errorf("Want Invalid Type Error, got: %v", err)
	}
}
