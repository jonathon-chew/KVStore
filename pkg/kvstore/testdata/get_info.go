package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strings"
)

func hashString(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func main() {
	content, err := os.ReadFile("/usr/share/dict/words")
	if err != nil {
		log.Printf("[ERROR]: Issues getting word list: %s", err.Error())
	}

	file_path := "./data.md"
	var binSize uint32 = 10

	if _, exists := os.Lstat(file_path); exists != nil {
		_, err := os.Create(file_path)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	var all_1s uint32 = 1
	for _, word := range strings.Split(string(content), "\n") {
		insert_word := strings.TrimSpace(word)
		bucket_number := hashString(insert_word) % binSize

		if bucket_number == all_1s && len(insert_word) > 4 {
			fmt.Println(insert_word)
		}
	}
}
