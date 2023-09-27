package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheMap map[string]CacheItem
	duration time.Duration
	mutex    sync.RWMutex
}

func NewCache(duration time.Duration) Cache {
	var cache = Cache{
		cacheMap: make(map[string]CacheItem),
		duration: duration,
	}

	go cache.reapLoop()

	return cache
}

func (ch *Cache) reapLoop() {
	for {
		ch.mutex.Lock()

		for index, item := range ch.cacheMap {
			res := time.Now().Sub(item.createdAt)

			if res > ch.duration {
				delete(ch.cacheMap, index)
			}
		}

		ch.mutex.Unlock()

		time.Sleep(ch.duration)
	}
}

func (ch *Cache) Get(key string) ([]byte, bool) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	val, ok := ch.cacheMap[key]
	return val.val, ok
}

func (ch *Cache) Add(key string, value []byte) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	ch.cacheMap[key] = CacheItem{
		createdAt: time.Now(),
		val:       value,
	}
}
