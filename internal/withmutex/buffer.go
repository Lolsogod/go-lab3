package withmutex

import (
	"awesomeProject/internal/tester"
	"sync"
)

type Buffer struct {
	data      interface{}
	isFull    bool
	mu        sync.Mutex
	writeCond *sync.Cond
	readCond  *sync.Cond
	wg        *sync.WaitGroup
}

func NewBuffer() *Buffer {
	b := &Buffer{
		isFull: false,
	}
	b.writeCond = sync.NewCond(&b.mu)
	b.readCond = sync.NewCond(&b.mu)
	return b
}

func (b *Buffer) Write(data interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for b.isFull {
		b.writeCond.Wait()
	}

	b.data = data
	b.isFull = true
	b.readCond.Signal()
}

func (b *Buffer) Read() interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	for !b.isFull {
		b.readCond.Wait()
	}

	data := b.data
	b.isFull = false
	b.writeCond.Signal()
	if b.wg != nil {
		b.wg.Done()
	}
	return data
}

func Test(messageCount int, messageSize int, threads int) {
	var wg sync.WaitGroup
	wg.Add(messageCount)
	buffer := NewBuffer()
	buffer.wg = &wg
	tester.RunBufferTest(tester.TestConfig{
		Name:         "Mutex Buffer",
		Buffer:       buffer,
		MessageCount: messageCount,
		MessageSize:  messageSize,
		ThreadCount:  threads,
		WaitGroup:    &wg,
	})
}
