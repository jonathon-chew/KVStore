package kvstore

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Entry struct {
	Key   string
	Value any
}

type HashTable struct {
	Buckets [][]Entry
	Size    int // How many elements in a single bucket before resizing, to keep the arrays short for faster look up
	mu      sync.RWMutex
}

func hashString(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (m *HashTable) bucketFor(key string) uint32 {
	hash := hashString(key)
	return hash % uint32(len(m.Buckets))
}

func (m *HashTable) exists(key string) (uint32, bool) {
	bucket_number := m.bucketFor(key)
	bucket := m.Buckets[bucket_number]

	for _, value := range bucket {
		if value.Key == key {
			return bucket_number, true
		}
	}

	return bucket_number, false
}

func (m *HashTable) resize() {
	oldBuckets := m.Buckets
	m.Buckets = make([][]Entry, len(m.Buckets)*2)

	for _, bucket_content := range oldBuckets {
		if len(bucket_content) == 0 {
			continue
		}

		for _, entry := range bucket_content {
			bucketNumber := m.bucketFor(entry.Key)

			m.Buckets[bucketNumber] = append(
				m.Buckets[bucketNumber],
				entry,
			)
		}
	}

}

func (m *HashTable) Add(key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	bucket_number, exists := m.exists(key)
	if exists {
		return fmt.Errorf("%s already exists", key)
	}

	bucket := m.Buckets[bucket_number]
	bucket = append(bucket, Entry{
		Key:   key,
		Value: value,
	})

	m.Buckets[bucket_number] = bucket
	m.Size++

	if float64(m.Size)/float64(len(m.Buckets)) >= 0.75 {
		m.resize()
	}

	return nil
}

func (m *HashTable) Get(key string) (any, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	bucket_number := m.bucketFor(key)
	bucket := m.Buckets[bucket_number]

	for _, value := range bucket {
		if value.Key == key {
			return value.Value, true
		}
	}

	return nil, false
}

func NewKVStore(size int) *HashTable {

	var kv HashTable

	kv.Buckets = make([][]Entry, size)
	kv.Size = size

	return &kv
}

func (m *HashTable) Remove(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	bucket_number, exists := m.exists(key)
	if !exists {
		return fmt.Errorf("%s does not exist", key)
	}

	bucket := m.Buckets[bucket_number]

	for idx, value := range bucket {
		if value.Key == key {
			bucket = append(bucket[:idx], bucket[idx+1:]...)
			m.Buckets[bucket_number] = bucket
			return nil
		}
	}

	m.Size--

	return nil
}

func (m *HashTable) Update(key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	buckerNumber, exists := m.exists(key)
	if !exists {
		m.Add(key, value)
		return nil
	}

	bucket := m.Buckets[buckerNumber]

	for idx, value := range bucket {
		if value.Key == key {
			bucket = append(bucket[:idx], bucket[idx+1:]...)
			m.Buckets[buckerNumber] = bucket
		}
	}

	bucket = append(bucket, Entry{
		Key:   key,
		Value: value,
	})

	m.Buckets[buckerNumber] = bucket

	return nil
}
