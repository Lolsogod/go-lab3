package withwg

import (
	"awesomeProject/internal/tester"
	"sync"
)

type Buffer struct {
	data   interface{}
	isFull bool
	mu     sync.Mutex
	wg     *sync.WaitGroup
}

func NewBuffer(wg *sync.WaitGroup) *Buffer {
	return &Buffer{
		wg:     wg,
		isFull: false,
	}
}

func (b *Buffer) Write(data interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for b.isFull {
		b.mu.Unlock()
		b.mu.Lock()
	}

	b.data = data
	b.isFull = true
}

func (b *Buffer) Read() interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	for !b.isFull {
		b.mu.Unlock()
		b.mu.Lock()
	}

	data := b.data
	b.isFull = false
	b.wg.Done()
	return data
}

func Test(messageCount int, messageSize int, threads int) {
	var wg sync.WaitGroup
	wg.Add(messageCount)
	buffer := NewBuffer(&wg)
	tester.RunBufferTest(tester.TestConfig{
		Name:         "WaitGroup Buffer",
		Buffer:       buffer,
		MessageCount: messageCount,
		MessageSize:  messageSize,
		ThreadCount:  threads,
		WaitGroup:    &wg,
	})
}
