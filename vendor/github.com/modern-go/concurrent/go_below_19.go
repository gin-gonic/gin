//+build !go1.9

package concurrent

import "sync"

type Map struct {
	lock sync.RWMutex
	data map[interface{}]interface{}
}

func NewMap() *Map {
	return &Map{
		data: make(map[interface{}]interface{}, 32),
	}
}

func (m *Map) Load(key interface{}) (elem interface{}, found bool) {
	m.lock.RLock()
	elem, found = m.data[key]
	m.lock.RUnlock()
	return
}

func (m *Map) Store(key interface{}, elem interface{}) {
	m.lock.Lock()
	m.data[key] = elem
	m.lock.Unlock()
}

