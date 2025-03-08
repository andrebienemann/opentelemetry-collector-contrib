package cache

type Cache struct {
	data map[string]*Record
}

func NewCache() Cache {
	data := make(map[string]*Record)
	return Cache{data}
}

func (cache *Cache) GetRecord(key string) (*Record, bool) {
	entry, ok := cache.data[key]
	return entry, ok
}

func (cache *Cache) PutRecord(key string, record *Record) {
	cache.data[key] = record
}
