package captcha

import (
	"sync"
	"time"
)

type StoreSyncMap struct {
	liveTime time.Duration
	m        *sync.Map
}

func NewStoreSyncMap(liveTime time.Duration) *StoreSyncMap {
	return &StoreSyncMap{liveTime: liveTime, m: new(sync.Map)}
}

type smv struct {
	t     time.Time
	Value string
}

func newSmv(v string) *smv {
	return &smv{t: time.Now(), Value: v}
}

func (s StoreSyncMap) rmExpire() {
	expireTime := time.Now().Add(-s.liveTime)
	s.m.Range(func(key, value interface{}) bool {
		if sv, ok := value.(*smv); ok && sv.t.Before(expireTime) {
			s.m.Delete(key)
		}
		return true
	})
}

func (s StoreSyncMap) Set(id string, value string) {
	s.rmExpire()
	s.m.Store(id, newSmv(value))
}

func (s StoreSyncMap) Get(id string, clear bool) string {
	v, ok := s.m.Load(id)
	if !ok {
		return ""
	}
	s.m.Delete(id)
	if sv, ok := v.(*smv); ok {
		return sv.Value
	}
	return ""
}

func (s StoreSyncMap) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}
