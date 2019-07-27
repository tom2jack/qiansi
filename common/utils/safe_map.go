package utils

import "sync"

type SafeStringMap struct {
	sync.RWMutex
	Map map[string]string
}

func NewSafeStringMap() *SafeStringMap {
	st := new(SafeStringMap)
	st.Map = make(map[string]string)
	return st
}

func (st *SafeStringMap) GET(key string) string {
	st.RLock()
	value := st.Map[key]
	st.RUnlock()
	return value
}

func (st *SafeStringMap) SET(key string, value string) {
	st.Lock()
	st.Map[key] = value
	st.Unlock()
}

func (st *SafeStringMap) SETNX(key string, value string) bool {
	st.Lock()
	if _, ok := st.Map[key]; ok {
		st.Unlock()
		return false
	}
	st.Map[key] = value
	st.Unlock()
	return true
}

func (st *SafeStringMap) DEL(key string) {
	st.Lock()
	delete(st.Map, key)
	st.Unlock()
}
