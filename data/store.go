package data

import (
	"audiobook-ingest/enum"
	"iter"
	"sync"
)

type StoreStruct struct {
	store map[string]*Item
	mutex sync.RWMutex
}

func NewStore() StoreStruct {
	return StoreStruct{
		store: make(map[string]*Item),
	}
}

func (s *StoreStruct) Add(item *Item) {
	id := item.ID()
	s.mutex.Lock()

	for _, existing := range s.store {
		if existing.FullPathName() == item.FullPathName() {
			s.mutex.Unlock()
			return
		}
	}
	s.store[id] = item
	s.mutex.Unlock()
	DetectingQueue <- id
}

func (s *StoreStruct) Remove(id string) (item *Item, exists bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	item, exists = s.store[id]
	if exists {
		delete(s.store, id)
		if SendMessage != nil {
			SendMessage(id, enum.RemoveItem, nil)
		}
	}
	return item, exists

}

func (s *StoreStruct) GetItem(id string) (item *Item, exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, exists = s.store[id]
	return
}

func (s *StoreStruct) ItemExists(fullPathName string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, item := range s.store {
		if item.FullPathName() == fullPathName {
			return true
		}
	}
	return false
}

func (s *StoreStruct) Count() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.store)
}

func (s *StoreStruct) Iterator() iter.Seq2[string, *Item] {
	return func(yield func(string, *Item) bool) {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
		for k, v := range s.store {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (s *StoreStruct) IsPathTracked(fullPathName string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, item := range s.store {
		if item.FullPathName() == fullPathName {
			return true
		}
	}
	return false
}
