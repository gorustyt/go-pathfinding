package path_finding

import (
	"sync"
	"time"
)

type PathPoint struct {
	X, Y int
}

type TracePoint struct {
	PathPoint
	IsJumpPoint bool
}

type DebugTrace interface {
	SetPathHandle(fn func(v *TracePoint))
	TracePath(x, y int, isJumpPoint bool) bool
	Wait()
	Exit()
	Start()
	Pause()
}
type debugTrace struct {
	closeCh  chan struct{}
	ch       chan *TracePoint
	pause    bool
	exit     func()
	interval time.Duration
	fn       func(v *TracePoint)

	wg sync.WaitGroup
}

func NewDebugTrace() DebugTrace {
	t := &debugTrace{
		interval: 10 * time.Millisecond,
		closeCh:  make(chan struct{}),
		ch:       make(chan *TracePoint)}
	t.exit = sync.OnceFunc(func() {
		close(t.closeCh)
	})
	t.wg.Add(1)
	go t.run()
	return t
}

func (trace *debugTrace) SetPathHandle(fn func(v *TracePoint)) {
	trace.fn = fn
}

func (trace *debugTrace) TracePath(x, y int, isJumpPoint bool) bool {
	select {
	case <-trace.closeCh:
		return false
	default:

	}
	select {
	case <-trace.closeCh:
		return false
	case trace.ch <- &TracePoint{
		IsJumpPoint: isJumpPoint,
		PathPoint: PathPoint{
			X: x,
			Y: y,
		}}:
	}
	return true
}

func (trace *debugTrace) Wait() {
	trace.wg.Wait()
}
func (trace *debugTrace) Exit() {
	trace.exit()
}
func (trace *debugTrace) Start() {
	trace.pause = false
}
func (trace *debugTrace) Pause() {
	trace.pause = true
}

func (trace *debugTrace) run() {
	timer := time.NewTicker(trace.interval)
	defer func() {
		timer.Stop()
		trace.wg.Done()
	}()
	for {
		select {
		case <-trace.closeCh:
			close(trace.ch)
			return
		case <-timer.C:
			if trace.pause {
				continue
			}
			select {
			case v, ok := <-trace.ch:
				if !ok {
					return
				}
				if trace.fn != nil {
					trace.fn(v)
				}
			default:

			}
		}

	}
}
