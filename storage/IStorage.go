package storage

type IStorage interface {
	AddItem(key string, val interface{})error
	RemoveItem(key string)
	GetItem(key string)(interface{}, error)
	GetAllItems()map[string] interface{}
}
