package kvstore

import (
	"testing"
)

func Test_AddKey(t *testing.T) {

	tests := []struct {
		key        string
		value      string
		shouldWork bool
	}{
		{key: "Example", value: "Hello World", shouldWork: true},
		{key: "Example", value: "Hello World", shouldWork: false},
	}

	var KVStore HashTable
	KVStore.Buckets = make([][]Entry, 100)

	for _, test := range tests {
		err := KVStore.Add(test.key, test.value)
		if err == nil && test.shouldWork == true {
			t.Logf("%s was successfully addded", test.key)
		} else if err != nil && test.shouldWork == true {
			t.Logf("%s errored but should have worked", test.key)
		} else if err != nil && test.shouldWork == false {
			t.Logf("%s was meant to error and did", test.key)
		} else if err == nil && test.shouldWork == false {
			t.Logf("%s was meant to error and failed to error", test.key)
		}
	}
}

func Test_MainFunctionality(t *testing.T) {
	var KVStore HashTable

	KVStore.Buckets = make([][]Entry, 10)

	key := "Example"
	value := "Hello World"

	KVStore.Add(key, value)
	entry, found := KVStore.Get(key)

	if found {
		t.Logf("Example was found with the value: %v\n", entry)
	} else {
		t.Errorf("Example was not found %v\n", entry)
		return
	}

	_, exists := KVStore.exists(key)

	if exists && found {
		t.Logf("Key was added and exists and can be found with the exist method")
	} else {
		t.Errorf("Key should have been added however it can't be found with the exist method")
	}

	KVStore.Remove(key)

	// Should no longer exist after being removed
	_, exists = KVStore.exists(key)
	if exists {
		t.Errorf("%s exists: %v\n", key, exists)
	} else {
		t.Logf("%s exists: %v\n", key, exists)
	}
}
