package resource

import (
	"context"
	"time"
)

type BeatStatus int

// Export
const (
	BeatInit BeatStatus = iota
	BeatRun
	BeatStop
)

type Tick struct {
	seq      int // 2147483648
	CreateTS time.Time
}

type CallbackFunc func(t Tick)

type Beat struct {
	Name     string
	CreateTS time.Time
	Status   BeatStatus
	maxSeq   int
	//
	context    context.Context
	interval   int // second
	tickChan   chan Tick
	stopped    chan int
	cancelFunc context.CancelFunc
	Callback   CallbackFunc
	lastTick   Tick
}

func (b *Beat) increase() {
	b.maxSeq++
}

// Export,
func NewBeatWithInterval(name string, interval int) *Beat {
	if interval <= 0 {
		panic("Beat interval must greater than 0")
	}
	ctx, cancalFunc := context.WithCancel(context.Background())
	return &Beat{
		Name:       name,
		context:    ctx,
		cancelFunc: cancalFunc,
		tickChan:   make(chan Tick),
		CreateTS:   time.Now(),
		interval:   interval,
		// Callback:   cb,
		stopped: make(chan int, 1),
	}
}

// Export
func (b *Beat) Tick() Tick {
	go b.run()
	for {
		select {
		case <-b.stopped:
			// stopped
			break
		case t := <-b.tickChan:
			b.lastTick = t
			b.Callback(t)
		}
	}
}

func (b *Beat) run() {
	if b.Status == BeatStop || b.Status == BeatRun {
		return
	}
	if b.Callback == nil {
		return
	}
	go func() {
		defer func() { b.stopped <- 1 }()
		tempTick := time.NewTicker(time.Duration(b.interval) * time.Second)
		b.Status = BeatRun
		for {
			select {
			case <-b.context.Done():
				// stop
				b.Status = BeatStop
				close(b.tickChan)
				break
			default:
				// do heartbeat
				t := <-tempTick.C
				b.tickChan <- Tick{seq: b.maxSeq, CreateTS: t}
				b.increase()
			}
		}
	}()
}

// Export, 阻塞
func (b *Beat) stop() {
	b.cancelFunc()
	// 等待Run退出
	<-b.stopped
}
