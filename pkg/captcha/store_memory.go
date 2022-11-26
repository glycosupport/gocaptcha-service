package captcha

import (
	"container/list"
	"sync"
	"time"
)

type idByTimeValue struct {
	timestamp time.Time
	id        string
}

type memoryStore struct {
	sync.RWMutex
	digitsById map[string]string
	idByTime   *list.List
	numStored  int
	collectNum int
	expiration time.Duration
}

func NewMemoryStore(collectNum int, expiration time.Duration) Store {
	s := new(memoryStore)
	s.digitsById = make(map[string]string)
	s.idByTime = list.New()
	s.collectNum = collectNum
	s.expiration = expiration
	return s
}

func (s *memoryStore) Set(id string, value string) error {
	s.Lock()
	s.digitsById[id] = value
	s.idByTime.PushBack(idByTimeValue{time.Now(), id})
	s.numStored++
	s.Unlock()
	if s.numStored > s.collectNum {
		go s.collect()
	}
	return nil
}

func (s *memoryStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == answer
}

func (s *memoryStore) Get(id string, clear bool) (value string) {
	if !clear {
		s.RLock()
		defer s.RUnlock()
	} else {
		s.Lock()
		defer s.Unlock()
	}
	value, ok := s.digitsById[id]
	if !ok {
		return
	}
	if clear {
		delete(s.digitsById, id)
	}
	return
}

func (s *memoryStore) collect() {
	now := time.Now()
	s.Lock()
	defer s.Unlock()
	for e := s.idByTime.Front(); e != nil; {
		e = s.collectOne(e, now)
	}
}

func (s *memoryStore) collectOne(e *list.Element, specifyTime time.Time) *list.Element {

	ev, ok := e.Value.(idByTimeValue)
	if !ok {
		return nil
	}

	if ev.timestamp.Add(s.expiration).Before(specifyTime) {
		delete(s.digitsById, ev.id)
		next := e.Next()
		s.idByTime.Remove(e)
		s.numStored--
		return next
	}
	return nil
}
