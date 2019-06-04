package utils

import "sync"

type SafeMap struct {
	sync.RWMutex
	Map map[string]string
}

func NewSafeMap() *SafeMap {
	st := new(SafeMap)
	st.Map = make(map[string]string)
	return st
}

func (st *SafeMap) GET(key string) string {
	st.RLock()
	value := st.Map[key]
	st.RUnlock()
	return value
}

func (st *SafeMap) SET(key string, value string) {
	st.Lock()
	st.Map[key] = value
	st.Unlock()
}

func (st *SafeMap) DEL(key string) {
	st.Lock()
	delete(st.Map, key)
	st.Unlock()
}