package perf

import (
	"fmt"
	"strconv"

	"github.com/jonathon-chew/KVStore/pkg/kvstore"
)

func NewPerf() *kvstore.HashTable {
	kv := kvstore.NewKVStore(1)

	return kv
}

func Simple(command, key, value string, kv *kvstore.HashTable) (any, error) {
	return_value, exists := kv.Get(key)
	if !exists {
		if err := kv.Add(key, 1); err != nil {
			return nil, fmt.Errorf("[ERROR]: key not found and unable to set to 0\n")
		}
		// fmt.Println("Added the key and initalised to 0")
		return "Added the key and initalised to 0", nil
	}

	var incr_value int = 1
	var err error
	if command == "INCRBY" && value != "" {
		incr_value, err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("[ERROR]: unable to incriment %s by %s\n", key, value)
		}
	}

	switch v := return_value.(type) {
	case int:
		if err := kv.Update(key, v+incr_value); err != nil {
			return nil, fmt.Errorf("[ERROR]: Unable to incriment key\n")
		}

		return nil, nil
	case string:
		number, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("[ERROR]: Unable to convert string to int, %s\n", v)
		}
		if remove_err := kv.Remove(key); remove_err != nil {
			return nil, fmt.Errorf("[ERROR]: Unable to remove key to re-add it with it's new value!, %s\n", v)
		}
		key_update_err := kv.Update(key, number+incr_value)
		if key_update_err != nil {
			if err := kv.Add(key, number); err != nil {
				return nil, fmt.Errorf("[ERROR]: Failed to add the new value, old one was removed and could not be added back\n")
			}
			return nil, fmt.Errorf("[ERROR]: Unable to incriment key\n")
		}
	default:
		return nil, fmt.Errorf("[ERROR]: Value is neither string which can be converted to an int, nor an int\n")
	}

	return nil, nil
}
