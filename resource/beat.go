package resource

import (
	"context"
	"time"
)

type BeatStatus int

const (
	BeatInit BeatStatus = iota
	BeatRun
	BeatStop
)

type Tick struct {
	seq      int // 2147483648
	CreateTS time.Time
}

type Beat struct {
	Name     string
	CreateTS time.Time
	Status   BeatStatus
	maxSeq   int
	Context  context.Context
	interval int // second
	tickChan chan Tick
}

func (b *Beat) Increase() {
	b.maxSeq++
}

func NewBeat(name string) *Beat {
	return &Beat{
		Name:    name,
		Context: context.Background(),
	}
}

func (b *Beat) Tick() Tick {
	return <-b.tickChan
}

func (b *Beat) Run() {
	if b.Status == BeatStop || b.Status == BeatRun {
		return
	}
	go func() {
		tempTick := time.NewTicker(time.Duration(b.interval) * time.Second)
		b.Status = BeatRun
		for {
			select {
			case <-b.Context.Done():
				// stop
				b.Status = BeatStop
				return
			default:
				// do heartbeat
				t := <-tempTick.C
				b.tickChan <- Tick{seq: b.maxSeq, CreateTS: t}
				b.Increase()
			}
		}
	}()

}
