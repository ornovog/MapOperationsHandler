package storage

import (
	"fmt"
	"sync"
)

type NavigationMap struct {
	mu 				sync.RWMutex
	keyToValueNode	map[string]*node
	values 			doublyLinkedList
}

func MakeNavigationMap() IStorage{
	navigationMap := NavigationMap{}
	navigationMap.values = *MakeDoublyList()
	return &navigationMap
}

func (n *NavigationMap) AddItem(key string, val string) error{
	n.mu.Lock()
	defer n.mu.Unlock()

	_, ok := n.keyToValueNode[key]
	if ok{
		return fmt.Errorf("key %s is already in the map", key)
	}

	n.values.AddEndNodeDLL(val)
	return nil
}

func (n *NavigationMap) RemoveItem(key string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	node, ok := n.keyToValueNode[key]
	if !ok{
		delete(n.keyToValueNode, key)
		n.values.RemoveNode(node)
	}
}

func (n *NavigationMap) GetItem(key string) (string, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	node, ok := n.keyToValueNode[key]
	if ok{
		return "", fmt.Errorf("key %s doesn't exist", key)
	}

	return node.data, nil
}

func (n *NavigationMap) GetAllItemsByOrder() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	values := make([]string, 0)
	for current := n.values.head ; current != nil; current = current.next{
		values = append(values, current.data)
	}
	return values
}