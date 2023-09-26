package main

import (
	"time"
)

type CacheItem struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheMap map[string]CacheItem
	duration time.Duration
}

func NewCache(duration time.Duration) Cache {
	return Cache{
		cacheMap: make(map[string]CacheItem),
		duration: duration,
	}
}

func (ch *Cache) Get(key string) ([]byte, bool) {
	val, ok := ch.cacheMap[key]
	return val.val, ok
}

func (ch *Cache) Add(key string, value []byte) {
	ch.cacheMap[key] = CacheItem{
		createdAt: time.Now(),
		val:       value,
	}
}
