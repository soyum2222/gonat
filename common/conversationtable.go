package common

import (
	_interface "gonat/interface"
	"sync"
)

type ConversationTable struct {
	l sync.RWMutex
	m map[uint32]_interface.Conversation
}

func (uc *ConversationTable) Init() {
	uc.m = map[uint32]_interface.Conversation{}
}

func (uc *ConversationTable) Store(key uint32, v _interface.Conversation) {

	uc.l.Lock()
	defer uc.l.Unlock()
	uc.m[key] = v

}

func (uc *ConversationTable) Load(key uint32) (_interface.Conversation, bool) {

	uc.l.RLock()
	defer uc.l.RUnlock()

	v, ok := uc.m[key]
	return v, ok
}

func (uc *ConversationTable) Delete(key uint32) {

	uc.l.Lock()
	defer uc.l.Unlock()
	delete(uc.m, key)
}

func (uc *ConversationTable) Range(f func(key uint32, value _interface.Conversation)) {

	for k, v := range uc.m {
		f(k, v)
	}
}
