/*
@author '彼时思默'
@time 2020/5/12 下午5:20
@describe:
*/
package entity

import "sync"

type PubMap struct {
	Body sync.Map
}

func (m *PubMap) Store(key string, sub *Pub) {
	m.Body.Store(key, sub)
}
func (m *PubMap) Load(key string) *Pub {
	if res, ok := m.Body.Load(key); ok {
		return res.(*Pub)
	}
	return nil
}
func (m *PubMap) Delete(key string) {
	m.Body.Delete(key)
}

func (m *PubMap) Range(f func(key string, sub *Pub)) {
	m.Body.Range(func(key, value interface{}) bool {
		f(key.(string), value.(*Pub))
		return true
	})
}

type SubMap struct {
	Body sync.Map
}

func (m *SubMap) Store(key string, sub *Sub) {
	m.Body.Store(key, sub)
}
func (m *SubMap) Load(key string) *Sub {
	if res, ok := m.Body.Load(key); ok {
		return res.(*Sub)
	}
	return nil
}
func (m *SubMap) Delete(key string) {
	m.Body.Delete(key)
}
func (m *SubMap) Range(f func(key string, sub *Sub)) {
	m.Body.Range(func(key, value interface{}) bool {
		f(key.(string), value.(*Sub))
		return true
	})
}
