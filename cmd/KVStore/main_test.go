package main

import (
	"testing"
)

func Test_Add(t *testing.T) {

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
