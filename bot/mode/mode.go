package mode

import "sync"

var Storage = struct {
	sync.RWMutex
	M map[int64]int
}{M: make(map[int64]int)}


const (
	InputAmount = iota
	WaitingForCategory 
	WaitingForPayment
	WaitingForConfirmation
)