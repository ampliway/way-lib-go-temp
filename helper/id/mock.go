package id

import (
	"sync"
)

var _ ID = (*Mock)(nil)

type Mock struct {
	mux           sync.Mutex
	sequenceCount int
	sequences     map[int]string
}

func NewMock() *Mock {
	return &Mock{
		mux:           sync.Mutex{},
		sequenceCount: -1,
		sequences:     map[int]string{},
	}
}

func (m *Mock) Random() string {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.sequenceCount++

	return m.sequences[m.sequenceCount]
}

func (m *Mock) ExpectRandom(s string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.sequences[len(m.sequences)] = s
}
