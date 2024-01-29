package path_finding

import (
	"sync"
	"time"
)

type DebugTrace struct {
	ch       chan [2]float32
	pause    bool
	exit     func()
	paths    [][2]float32
	interval time.Duration
}

func NewDebugTrace() *DebugTrace {
	t := &DebugTrace{
		interval: 100 * time.Millisecond,
		ch:       make(chan [2]float32)}
	t.exit = sync.OnceFunc(func() {
		close(t.ch)
	})
	return t
}

func (trace *DebugTrace) Send(pos [2]float32) {
	trace.ch <- pos
}

func (trace *DebugTrace) GetPaths() [][2]float32 {
	return trace.paths
}

func (trace *DebugTrace) Exit() {
	trace.exit()
	trace.paths = nil
}
func (trace *DebugTrace) Start() {
	trace.pause = false
}
func (trace *DebugTrace) Pause() {
	trace.pause = true
}

func (trace *DebugTrace) run() {
	timer := time.NewTimer(trace.interval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if trace.pause {
				continue
			}
			select {
			case v, ok := <-trace.ch:
				if !ok {
					return
				}
				trace.paths = append(trace.paths, v)
			}
		}

	}
}
