package kvstore

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

type Store interface {
	Add(key string, value any) error
	Get(key string) (any, bool)
	Remove(key string) error
	Update(key string, value any) error
}

type entry struct {
	key       string
	value     any
	expiresAt time.Time
}

type hashTable struct {
	buckets [][]entry
	size    int // How many elements in a single bucket before resizing, to keep the arrays short for faster look up
	mu      sync.RWMutex
}

// hashString returns the FNV-1a hash for key.
func hashString(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// bucketFor returns the bucket index for key.
func (m *hashTable) bucketFor(key string) uint32 {
	hash := hashString(key)
	return hash % uint32(len(m.buckets))
}

// exists reports whether key is present and returns its bucket index.
func (m *hashTable) exists(key string) (uint32, bool) {
	bucket_number := m.bucketFor(key)
	bucket := m.buckets[bucket_number]

	for _, value := range bucket {
		if value.key == key {
			return bucket_number, true
		}
	}

	return bucket_number, false
}

// remove deletes key from the store without acquiring the mutex.
func (m *hashTable) remove(key string) error {
	bucket_number, exists := m.exists(key)
	if !exists {
		return fmt.Errorf("%s does not exist", key)
	}

	bucket := m.buckets[bucket_number]

	for idx, value := range bucket {
		if value.key == key {
			bucket = append(bucket[:idx], bucket[idx+1:]...)
			m.buckets[bucket_number] = bucket
			m.size--
			return nil
		}
	}

	return nil
}

// resize doubles the bucket count and rehashes all entries.
func (m *hashTable) resize() {
	oldBuckets := m.buckets
	m.buckets = make([][]entry, len(m.buckets)*2)

	for _, bucket_content := range oldBuckets {
		if len(bucket_content) == 0 {
			continue
		}

		for _, entry := range bucket_content {
			bucketNumber := m.bucketFor(entry.key)

			m.buckets[bucketNumber] = append(
				m.buckets[bucketNumber],
				entry,
			)
		}
	}

}

// Add inserts value under key and returns an error if key already exists.
func (m *hashTable) Add(key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	bucket_number, exists := m.exists(key)
	if exists {
		return fmt.Errorf("%s already exists", key)
	}

	bucket := m.buckets[bucket_number]
	bucket = append(bucket, entry{
		key:   key,
		value: value,
	})

	m.buckets[bucket_number] = bucket
	m.size++

	if float64(m.size)/float64(len(m.buckets)) >= 0.75 {
		m.resize()
	}

	return nil
}

// Get returns the value stored under key and whether the key exists.
func (m *hashTable) Get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	bucket_number := m.bucketFor(key)
	bucket := m.buckets[bucket_number]

	for _, value := range bucket {
		if value.key == key {
			return value.value, true
		}
	}

	return nil, false
}

// Exists reports whether key is present in the store.
func (m *hashTable) Exists(key string) bool {
	if _, exists := m.Get(key); !exists {
		return false
	}

	return true
}

// NewKVStore creates a key-value store with the requested initial bucket count.
func NewKVStore(size int) Store {

	var kv hashTable

	if size < 1 {
		size = 1
	}

	kv.buckets = make([][]entry, size)

	return &kv
}

// Remove deletes key from the store.
func (m *hashTable) Remove(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.remove(key)
}

// TTL returns the time-to-live information for key.
func (m *hashTable) TTL(key string) int64 {
	m.mu.RLocker()
	defer m.mu.RUnlock()

	buckernumber := m.bucketFor(key)
	bucket := m.buckets[buckernumber]

	for idx, entry := range bucket {
		if idx == int(buckernumber) {
			return entry.expiresAt.UTC().Unix()
		}
	}

	return 0
}

// Update replaces the value stored under key, or adds key if it is absent.
func (m *hashTable) Update(key string, value any) error {
	m.mu.Lock()

	buckerNumber, exists := m.exists(key)
	if !exists {
		m.mu.Unlock()
		m.Add(key, value)
		return nil
	}

	defer m.mu.Unlock()

	bucket := m.buckets[buckerNumber]

	for idx, value := range bucket {
		if value.key == key {
			bucket = append(bucket[:idx], bucket[idx+1:]...)
			m.buckets[buckerNumber] = bucket
		}
	}

	bucket = append(bucket, entry{
		key:   key,
		value: value,
	})

	m.buckets[buckerNumber] = bucket

	return nil
}
