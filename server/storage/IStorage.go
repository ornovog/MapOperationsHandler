package storage

type keyValuePair struct {
	Key string
	Value interface{}
}

type IStorage interface {
	AddItem(key string, val string)error
	RemoveItem(key string)
	GetItem(key string)(string, error)
	GetAllItemsByOrder()[]string
}
