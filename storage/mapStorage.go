package storage

import (
	"fmt"
	"github.com/orcaman/concurrent-map"
)

type mapStorage struct {
	cmap cmap.ConcurrentMap
}

func MakeMapStorage() IStorage{
	cmap := cmap.New()
	mapStorage := mapStorage{
		cmap: cmap,
	}
	return &mapStorage
}

func (ms *mapStorage)AddItem(key string, val interface{}) error{
	ok := ms.cmap.SetIfAbsent(key, val)
	if !ok{
		return fmt.Errorf("key %s already exists", key)
	}

	return nil
}

func (ms *mapStorage)RemoveItem(key string){
	ms.cmap.Remove(key)
}

func (ms *mapStorage)GetItem(key string)(interface{}, error){
	val, ok := ms.cmap.Get(key)
	if !ok {
		return "", fmt.Errorf("key %s doesn't exist", key)
	}

	return val.(string), nil
}

func (ms *mapStorage)GetAllItems()map[string]interface{}{
	return ms.cmap.Items()
}
