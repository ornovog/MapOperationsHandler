package storage

import (
	"fmt"
	"github.com/zhangyunhao116/skipmap"
)

type skipmapStorage struct {
	skipmap *skipmap.StringMap
}

func MakeSkipmapStorageStorage() IStorage{
	skipmap := skipmap.NewString()
	skipmapStorage := skipmapStorage{
		skipmap: skipmap,
	}
	return &skipmapStorage
}

func (ss *skipmapStorage)AddItem(key string, val interface{}) error{
	ss.skipmap.Store(key, val)
	return nil
}

func (ss *skipmapStorage)RemoveItem(key string){
	ss.skipmap.Delete(key)
}

func (ss *skipmapStorage)GetItem(key string)(interface{}, error){
	val, ok := ss.skipmap.Load(key)
	if !ok {
		return "", fmt.Errorf("key %s doesn't exist", key)
	}

	return val.(string), nil
}

func (ss *skipmapStorage)GetAllItemsByOrder()[]keyValuePair{
	allItems := make([]keyValuePair, 0)
	ss.skipmap.Range(func(key string, value interface{}) bool {
		keyValuePair := keyValuePair{Key: key, Value: value}
		allItems = append(allItems, keyValuePair)
		return true
	})
	return allItems
}
