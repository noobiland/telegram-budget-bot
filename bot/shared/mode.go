package shared

import "sync"

var Mode = struct {
	sync.RWMutex
	M map[int64]string
}{M: make(map[int64]string)}
