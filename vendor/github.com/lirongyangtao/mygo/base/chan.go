package base

import (
	"fmt"
	"sync"
)

// 无限缓存channel
type UnboundChan struct {
	ch     chan interface{}
	lock   *sync.Mutex
	buffer []interface{}
}

func NewUnboundChan(size int32) *UnboundChan {
	if size <= 0 || size > 4096 {
		size = 4096
	}
	ch := &UnboundChan{
		ch:   make(chan interface{}, size),
		lock: &sync.Mutex{},
	}
	return ch
}

func (ch *UnboundChan) Put(element interface{}) {
	ch.flush()
	select {
	case ch.ch <- element:
		fmt.Println("投递ele", element)
		return
	default:
		ch.lock.Lock()
		ch.buffer = append(ch.buffer, element)
		ch.lock.Unlock()
	}
}

func (ch *UnboundChan) flush() {
	ch.lock.Lock()
	if len(ch.buffer) <= 0 {
		return
	}
	defer ch.lock.Unlock()
	for {
		select {
		case ch.ch <- ch.buffer[0]:
			ch.buffer = ch.buffer[1:]
			continue
		default:
			break
		}
	}
}

func (ch *UnboundChan) Get() <-chan interface{} {
	ch.flush()
	return ch.ch
}
