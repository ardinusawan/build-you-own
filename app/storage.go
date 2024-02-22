package main

type Storage interface {
	Set(key string, val string)
	Get(key string) string
}

type Pair struct {
	Key string
	Val string
}

type MemoryStorage struct {
	items []Pair
}

func (m *MemoryStorage) Set(key string, val string) {
	for i, item := range m.items {
		if item.Key == key {
			m.items[i].Val = val
			return
		}
	}
	m.items = append(m.items, Pair{key, val})
}

func (m *MemoryStorage) Get(key string) string {
	for _, item := range m.items {
		if item.Key == key {
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
