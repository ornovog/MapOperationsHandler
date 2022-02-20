package storage

import (
	"fmt"
	"github.com/orcaman/concurrent-map"
)

type cmapStorage struct {
	cmap cmap.ConcurrentMap
}

/*func MakeCmapStorage() IStorage{
	cmap := cmap.New()
	cmapStorage := cmapStorage{
		cmap: cmap,
	}
	return &cmapStorage
}*/

func (ms *cmapStorage)AddItem(key string, val interface{}) error{
	ok := ms.cmap.SetIfAbsent(key, val)
	if !ok{
		return fmt.Errorf("key %s already exists", key)
	}

	return nil
}

func (ms *cmapStorage)RemoveItem(key string){
	ms.cmap.Remove(key)
}

func (ms *cmapStorage)GetItem(key string)(interface{}, error){
	val, ok := ms.cmap.Get(key)
	if !ok {
		return "", fmt.Errorf("key %s doesn't exist", key)
	}

	return val.(string), nil
}

func (ms *cmapStorage)GetAllItemsByOrder()[]keyValuePair{
	allItems := make([]keyValuePair, 0)
	for key, val := range ms.cmap.Items() {
		keyValuePair := keyValuePair{Key: key, Value: val}
		allItems = append(allItems, keyValuePair)
	}

	return allItems
}
