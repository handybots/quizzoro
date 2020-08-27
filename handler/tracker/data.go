package tracker

import "sync"

var Data = dataCache{
	m: make(map[int64]string),
}

type dataCache struct {
	sync.RWMutex
	m map[int64]string
}

func (u *dataCache) Get(id int64) (string, bool) {
	u.RLock()
	defer u.RUnlock()
	data, ok := u.m[id]
	return data, ok
}

func (u *dataCache) Set(id int64, v string) {
	u.Lock()
	u.m[id] = v
	u.Unlock()
}

func (u *dataCache) Del(id int64) {
	u.Lock()
	delete(u.m, id)
	u.Unlock()
}
