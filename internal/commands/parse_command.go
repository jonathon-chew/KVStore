package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jonathon-chew/KVStore/internal/kvstore"
)

func check_len_min(check string, minimum_length int) bool {
	return len(check) >= minimum_length
}

func ParseCommand(request string, kv *kvstore.HashTable) (any, error) {

	var command, key, value string

	split_string := strings.Split(request, " ")
	fmt.Println(split_string)

	if len(split_string) >= 2 {
		return nil, fmt.Errorf("[ERROR]: not enough commands, minimum 2\n")
	}

	command = split_string[0]
	key = split_string[1]

	if len(split_string) == 3 {
		value = split_string[2]
	}

	switch strings.Split(command, " ")[0] {
	case "SET":
		fmt.Println("User would like to set a key")
		if !check_len_min(key, 1) || !check_len_min(value, 1) {
			return nil, fmt.Errorf("[ERROR]: not enough commands to SET, minimum 3\n")
		}
		if err := kv.Add(key, value); err != nil {
			return nil, err
		}
		return "Successfully added", nil
	case "GET":
		fmt.Println("User would like to get a key")
		value, err := kv.Get(key)
		if !err {
			return nil, fmt.Errorf("[ERROR]: No key for: %s, %v\n", key, err)
		}
		return value, nil
	case "GETDEL":
		fmt.Println("User would like to get a key and delete it")
		value, err := kv.Get(key)
		if !err {
			return nil, fmt.Errorf("[ERROR]: No key for: %s, %v\n", key, err)
		}
		kv.Remove(key)
		return value, nil
	case "GETSET":
		fmt.Println("User would like to get a key and delete it")
		return_value, err := kv.Get(key)
		if !err {
			return nil, fmt.Errorf("[ERROR]: No key for: %s, %v", key, err)
		}
		if !check_len_min(key, 1) || !check_len_min(value, 1) {
			return nil, fmt.Errorf("[ERROR]: not enough commands to SET, minimum 3\n")
		}
		kv.Add(key, value)
		return return_value, nil
	case "DEL":
		fmt.Println("User would like to del a key")
		if err := kv.Remove(key); err != nil {
			return nil, err
		}
		return "Successfully added", nil
	case "STRLEN":
		return_value, exists := kv.Get(key)
		if !exists {
			return nil, fmt.Errorf("[ERROR]: unable to get the value of: %s\n", key)
		}
		value, ok := return_value.(string)
		if ok {
			return len(value), nil
		}
		return nil, fmt.Errorf("[ERROR]: Unable to get the str length of %s as it's value %s is not a string\n", key, return_value)
	case "INCR":
		fmt.Println("User would like to incriment a value, or set it's value to 0")
		return_value, exists := kv.Get(key)
		if !exists {
			if err := kv.Add(key, 0); err != nil {
				return nil, fmt.Errorf("[ERROR]: key not found and unable to set to 0\n")
			}
		}

		switch v := return_value.(type) {
		case int:
			if err := kv.Add(key, v+1); err != nil {
				return nil, fmt.Errorf("[ERROR]: Unable to incriment key\n")
			}
		case string:
			number, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("[ERROR]: Unable to convert string to int, %s\n", v)
			}
			if remove_err := kv.Remove(key); remove_err != nil {
				return nil, fmt.Errorf("[ERROR]: Unable to remove key to re-add it with it's new value!, %s\n", v)
			}
			key_add_err := kv.Add(key, number)
			if key_add_err != nil {
				if err := kv.Add(key, value); err != nil {
					return nil, fmt.Errorf("[ERROR]: Failed to add the new value, old one was removed and could not be added back\n")
				}
				return nil, fmt.Errorf("[ERROR]: Unable to incriment key\n")
			}
		default:
			return nil, fmt.Errorf("[ERROR]: Value is neither string which can be converted to an int, nor an int\n")
		}
	}
	return "", nil
}
