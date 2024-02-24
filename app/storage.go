package main

import (
	"time"
)

type Storage interface {
	Set(key string, val string, px *int64)
	Get(key string) string
}

type Pair struct {
	Key      string
	Val      string
	ExpireIn *int64 // UnixMilli
}

type MemoryStorage struct {
	items []Pair
}

func (m *MemoryStorage) Set(key string, val string, px *int64) {
	var expireIn *int64
	if px != nil {
		expirePtr := time.Now().UnixMilli() + *px
		expireIn = &expirePtr
	}

	for i, item := range m.items {
		if item.Key == key {
			m.items[i].Val = val
			m.items[i].ExpireIn = expireIn
			return
		}
	}
	m.items = append(m.items, Pair{key, val, expireIn})
}

func (m *MemoryStorage) Get(key string) string {
	for _, item := range m.items {
		if item.Key == key && item.ExpireIn == nil {
			return item.Val
		}

		expireIn := item.ExpireIn
		if item.Key == key && *expireIn <= time.Now().UnixMilli() {
			return "-1"
		}

		if item.Key == key && *expireIn > time.Now().UnixMilli() {
			return item.Val
		}
	}
	return ""
}

func NewMemoryStorage() *MemoryStorage {
	m := MemoryStorage{
		items: []Pair{},
	}
	return &m
}
