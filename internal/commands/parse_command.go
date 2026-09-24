package commands

import (
	"fmt"
	"strings"

	"github.com/jonathon-chew/KVStore/internal/kvstore"
)

func ParseCommand(request string, kv *kvstore.HashTable) (any, error) {

	split_string := strings.Split(request, " ")
	fmt.Println(split_string)
	command := split_string[0]
	key := split_string[1]

	switch strings.Split(command, " ")[0] {
	case "SET":
		fmt.Println("User would like to set a key")
		err := kv.Add(split_string[1], split_string[2])
		if err != nil {
			return nil, err
		}

		return "Successfully added", nil
	case "GET":
		fmt.Println("User would like to get a key")
		value, err := kv.Get(key)
		if !err {
			return nil, fmt.Errorf("[ERROR]: No key for: %s, %v", key, err)
		}
		return value, nil
	case "DEL":
		fmt.Println("User would like to del a key")
		err := kv.Remove(key)
		if err != nil {
			return nil, err
		}

		return "Successfully added", nil
	}

	return "", nil
}
