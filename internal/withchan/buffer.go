package withchan

import (
	"awesomeProject/internal/tester"
	"sync"
)

type Buffer struct {
	ch chan interface{}
	wg *sync.WaitGroup
}

func NewBuffer() *Buffer {
	return &Buffer{
		ch: make(chan interface{}, 1),
	}
}

func (b *Buffer) Write(data interface{}) {
	b.ch <- data
}

func (b *Buffer) Read() interface{} {
	data := <-b.ch
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
		Name:         "Channel Buffer",
		Buffer:       buffer,
		MessageCount: messageCount,
		MessageSize:  messageSize,
		ThreadCount:  threads,
		WaitGroup:    &wg,
	})
}
