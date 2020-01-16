package conn

import (
	_interface "gonat/interface"
	"sync"
)

type userTable struct {
	l sync.RWMutex
	m map[uint32]_interface.Conversation
}

func (uc *userTable) Init() {
	uc.m = map[uint32]_interface.Conversation{}
}

func (uc *userTable) Store(key uint32, v _interface.Conversation) {

	uc.l.Lock()
	defer uc.l.Unlock()
	uc.m[key] = v

}

func (uc *userTable) Load(key uint32) (_interface.Conversation, bool) {

	uc.l.RLock()
	defer uc.l.RUnlock()

	v, ok := uc.m[key]
	return v, ok
}

func (uc *userTable) Delete(key uint32) {

	uc.l.Lock()
	defer uc.l.Unlock()
	delete(uc.m, key)
}

func (uc *userTable) Range(f func(key uint32, value _interface.Conversation)) {

	for k, v := range uc.m {
		f(k, v)
	}
}
