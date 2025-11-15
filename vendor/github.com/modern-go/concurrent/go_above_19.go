//+build go1.9

package concurrent

import "sync"

type Map struct {
	sync.Map
}

func NewMap() *Map {
	return &Map{}
}