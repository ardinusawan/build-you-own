package main

import (
	"time"
)

type Storage interface {
	Set(key string, val string, px *int64)
	Get(key string) string
}

type Data struct {
	Val      string
	ExpireIn *int64 // UnixMilli
}

type MemoryStorage struct {
	key map[string]Data
}

func (m *MemoryStorage) Set(key string, val string, px *int64) {
	var expireIn *int64
	if px != nil {
		expirePtr := time.Now().UnixMilli() + *px
		expireIn = &expirePtr
	}

	m.key[key] = Data{Val: val, ExpireIn: expireIn}
}

func (m *MemoryStorage) Get(key string) string {
	if data, ok := m.key[key]; ok && data.ExpireIn == nil {
		return data.Val
	}

	if data, ok := m.key[key]; ok && *data.ExpireIn <= time.Now().UnixMilli() {
		m.Del(key)
		return "-1"
	}

	if data, ok := m.key[key]; ok && *data.ExpireIn > time.Now().UnixMilli() {
		return data.Val
	}
	return ""
}

func (m *MemoryStorage) Del(key string) {
	delete(m.key, key)
}

func NewMemoryStorage() *MemoryStorage {
	m := MemoryStorage{
		key: make(map[string]Data),
	}
	return &m
}
