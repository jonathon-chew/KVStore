package kvstore

import (
	"log"
	"os"
	"strconv"
	"strings"
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

func Test_Resize(t *testing.T) {
	// 	Create a KVStore with 10 buckets.
	kvstore := NewKVStore(10)

	// Insert 8 unique keys.
	for i := range 8 {
		update := strconv.Itoa(i)
		kvstore.Add(update, update)
	}

	// Verify the bucket count has doubled.
	if len(kvstore.Buckets) != 20 {
		t.Fatalf("[ERROR]: Size is: %d not 20", len(kvstore.Buckets))
	}

	// Retrieve all 8 keys.
	for i := range 8 {
		update := strconv.Itoa(i)
		value, exists := kvstore.Get(update)

		if !exists {
			t.Fatalf("[ERROR]: %s doesn't exist", update)
		}

		// Verify every value is correct.
		if value != update {
			t.Errorf("[ERROR]: expected /%s/ got /%s/", update, value)
		}
	}

}

func findCollisions(binSize uint32, number_of_words int) ([]string, error) {
	content, err := os.ReadFile("/usr/share/dict/words")
	if err != nil {
		log.Printf("[ERROR]: Issues getting word list: %s", err.Error())
		return []string{}, err
	}

	var all_1s uint32 = 1
	var return_words []string
	for _, word := range strings.Split(string(content), "\n") {
		insert_word := strings.TrimSpace(word)
		bucket_number := hashString(insert_word) % binSize

		if bucket_number == all_1s && len(insert_word) > 4 {
			return_words = append(return_words, insert_word)
			if len(return_words) >= number_of_words {
				return return_words, nil
			}
		}
	}

	return return_words, nil
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func Test_Collisions(t *testing.T) {

	// Keep consistancy!
	var bucket_size int = 10

	// Find three different keys that hash to the same bucket.
	words, err := findCollisions(uint32(bucket_size), 4)
	if err != nil {
		t.Fatalf("[ERROR]: Unable to get words %s", err.Error())
	}

	// 	Create a KVStore with 10 buckets.
	kv := NewKVStore(bucket_size)

	for _, word := range words {
		// Insert all three keys with different values.
		kv.Add(word, reverse(word))
	}

	// Retrieve each key.
	for _, word := range words {
		value, exists := kv.Get(word)
		if !exists {
			t.Fatalf("[ERROR] %s does not exist", word)
		}

		// Verify each returns its own value.
		if value != reverse(word) {
			t.Fatalf("[ERROR]: Expected /%s/ Got /%s/", reverse(word), value)
		}
	}

	// Remove the middle key.
	err = kv.Remove(words[1])
	if err != nil {
		t.Fatalf("[ERROR]: Didn't remove the key properly")
	}

	if _, exists := kv.Get(words[1]); exists {
		t.Fatalf("[ERROR]: Removed key still exists")
	}

	// Verify the other two still exist.
	for _, idx := range []int{0, 2} {
		value, exists := kv.Get(words[idx])
		if !exists {
			t.Fatalf("[ERROR] /%s/ does not exist", words[idx])
		}

		// Verify each returns its own value.
		if value != reverse(words[idx]) {
			t.Fatalf("[ERROR]: Expected /%s/ Got /%s/", reverse(words[idx]), value)
		}
	}
}

func Test_Collision_Removal(t *testing.T) {}

func Test_Update(t *testing.T) {}

func Test_Remove(t *testing.T) {}

func Test_Resize_With_Collisions(t *testing.T) {}

func Test_Concurrent_Operations(t *testing.T) {}

func Fuzz_Differential(t *testing.F) {}

func FuzzAdding(f *testing.F) {
	f.Add("Hello", "greioi")
	f.Add("World", "ioierg")

	f.Fuzz(func(t *testing.T, key, value string) {
		kvstore := NewKVStore(10)

		err := kvstore.Add(key, value)
		if err != nil {
			t.Error("[ERROR]: ", err.Error())
			return
		}

		response, exists := kvstore.Get(key)
		if !exists {
			t.Error("[ERROR]: Key does not exist")
			return
		}

		if response != value {
			t.Errorf("[ERROR]: expected /%s/ got /%s/", value, response)
		}
	})
}
