package nosync

import (
	"awesomeProject/internal/tester"
)

type Buffer struct {
	data   interface{}
	isFull bool
}

func NewBuffer() *Buffer {
	return &Buffer{
		data:   nil,
		isFull: false,
	}
}

func (b *Buffer) Write(data interface{}) {
	b.data = data
	b.isFull = true
}

func (b *Buffer) Read() interface{} {
	b.isFull = false
	return b.data
}

func Test(messageCount, messageSize, threadCount int) int {
	buffer := NewBuffer()
	tester.RunBufferTest(tester.TestConfig{
		Name:         "NoSync Buffer",
		Buffer:       buffer,
		MessageCount: messageCount,
		MessageSize:  messageSize,
		ThreadCount:  threadCount,
	})

	successfulReads := 0
	
	for i := 0; i < threadCount; i++ {
		go func(id int) {
			messagesPerThread := messageCount / threadCount
			for j := 0; j < messagesPerThread; j++ {
				if buffer.Read() != nil {
					successfulReads++
				}
			}
		}(i)
	}
	
	return messageCount - successfulReads 
