package storage

type keyValuePair struct {
	Key string
	Value interface{}
}

type IStorage interface {
	AddItem(key string, val interface{})error
	RemoveItem(key string)
	GetItem(key string)(interface{}, error)
	GetAllItemsByOrder()[]keyValuePair
}
