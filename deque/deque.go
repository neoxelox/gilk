// Package deque forked from https://github.com/oleiade/lane
package deque

import (
	"container/list"
	"sync"
)

// Deque is a head-tail linked list data structure implementation.
// It is based on a double linked list container, so that every
// operation's time complexity is O(1).
//
// Every operation and Deque instances are synchronized and
// safe for concurrent usage.
type Deque struct {
	sync.RWMutex
	container *list.List
	capacity  int
}

// New creates a Deque
func New() *Deque {
	return NewCapped(-1)
}

// NewCapped creates a Deque with the specified capacity limit
func NewCapped(capacity int) *Deque {
	return &Deque{
		container: list.New(),
		capacity:  capacity,
	}
}

// Append inserts item at the back of the Deque in a O(1) time complexity,
// returning true if successful or false if the deque is at max capacity
func (s *Deque) Append(item interface{}) bool {
	s.Lock()
	defer s.Unlock()

	if s.capacity < 0 || s.container.Len() < s.capacity {
		s.container.PushBack(item)
		return true
	}

	return false
}

// Prepend inserts item at the front of the Deque in a O(1) time complexity,
// returning true if successful or false if the deque is at max capacity
func (s *Deque) Prepend(item interface{}) bool {
	s.Lock()
	defer s.Unlock()

	if s.capacity < 0 || s.container.Len() < s.capacity {
		s.container.PushFront(item)
		return true
	}

	return false
}

// IterFirst returns the front internal list iterator
func (s *Deque) IterFirst() *list.Element {
	s.RLock()
	defer s.RUnlock()

	return s.container.Front()
}

// IterLast returns the back internal list iterator
func (s *Deque) IterLast() *list.Element {
	s.RLock()
	defer s.RUnlock()

	return s.container.Back()
}

// Pop removes the last element of the Deque in a O(1) time complexity
func (s *Deque) Pop() interface{} {
	s.Lock()
	defer s.Unlock()

	var item interface{} = nil

	if lastContainerItem := s.container.Back(); lastContainerItem != nil {
		item = s.container.Remove(lastContainerItem)
	}

	return item
}

// Shift removes the first element of the Deque in a O(1) time complexity
func (s *Deque) Shift() interface{} {
	s.Lock()
	defer s.Unlock()

	var item interface{} = nil

	if firstContainerItem := s.container.Front(); firstContainerItem != nil {
		item = s.container.Remove(firstContainerItem)
	}

	return item
}

// First returns the first value stored in the Deque in a O(1) time complexity
func (s *Deque) First() interface{} {
	s.RLock()
	defer s.RUnlock()

	item := s.container.Front()
	if item != nil {
		return item.Value
	}
	return nil
}

// Last returns the last value stored in the Deque in a O(1) time complexity
func (s *Deque) Last() interface{} {
	s.RLock()
	defer s.RUnlock()

	item := s.container.Back()
	if item != nil {
		return item.Value
	}
	return nil
}

// Size returns the actual Deque size
func (s *Deque) Size() int {
	s.RLock()
	defer s.RUnlock()

	return s.container.Len()
}

// Capacity returns the capacity of the Deque, or -1 if unlimited
func (s *Deque) Capacity() int {
	s.RLock()
	defer s.RUnlock()
	return s.capacity
}

// Empty checks if the Deque is empty
func (s *Deque) Empty() bool {
	s.RLock()
	defer s.RUnlock()

	return s.container.Len() == 0
}

// Full checks if the Deque is full
func (s *Deque) Full() bool {
	s.RLock()
	defer s.RUnlock()

	return s.capacity >= 0 && s.container.Len() >= s.capacity
}
