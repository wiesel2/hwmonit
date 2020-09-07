package resource

import (
	"context"
	"sync"
	"time"
)

// var logger = common.GetLogger()

// BeatStatus Export
type BeatStatus int

// Export
const (
	BeatInit BeatStatus = iota
	BeatRun
	BeatStop
)

// Tick epxort
type Tick struct {
	seq      int // 2147483648
	CreateTS time.Time
}

// CallbackFunc Export
// Tick Callback
type CallbackFunc func(t Tick)

//Beat Export
type Beat struct {
	Name       string
	CreateTS   time.Time
	status     BeatStatus
	maxSeq     int
	context    context.Context
	interval   int // second
	tickChan   chan *Tick
	cancelFunc context.CancelFunc
	Callback   CallbackFunc
	lastTick   Tick
	ticker     *time.Ticker
	wg         sync.WaitGroup
}

func (b *Beat) increase() {
	b.maxSeq++
}

// NewBeatWithInterval Export
// Creator
func NewBeatWithInterval(name string, interval int, cb CallbackFunc) *Beat {
	if interval <= 0 {
		panic("Beat interval must greater than 0")
	}
	ctx, cancalFunc := context.WithCancel(context.Background())
	return &Beat{
		Name:       name,
		context:    ctx,
		cancelFunc: cancalFunc,
		tickChan:   make(chan *Tick, 1),
		CreateTS:   time.Now(),
		interval:   interval,
		Callback:   cb,
		wg:         sync.WaitGroup{},
	}
}

func (b *Beat) addOneTick(t time.Time) {
	b.tickChan <- &Tick{seq: b.maxSeq, CreateTS: t}
	b.increase()
}

// Func of beat
func (b *Beat) tick() {
	b.wg.Add(1)
	for {
		select {
		case <-b.context.Done():
			// stop
			b.status = BeatStop
			// logger.Debugf("Beat tick of %s stopped", b.Name)
			b.wg.Done()
			return
		case t := <-b.ticker.C:
			// do heartbeat
			// logger.Debugf("Beat tick of %s added", b.Name)
			b.addOneTick(t)
		}
	}
}

func (b *Beat) listen() {
	b.wg.Add(1)
	for {
		select {
		case <-b.context.Done():
			// logger.Debugf("Beat listen of %s stopped", b.Name)
			b.wg.Done()
			return
		case t, ok := <-b.tickChan:
			if ok == true {
				b.lastTick = *t
				// logger.Debugf("Beat listen of %s got one tick", b.Name)
				b.Callback(*t)
			}

		}
	}
}

// Start export, non-blocking
func (b *Beat) Start() {
	if b.status == BeatStop || b.status == BeatRun {
		return
	}
	if b.Callback == nil {
		return
	}
	b.ticker = time.NewTicker(time.Duration(b.interval) * time.Second)
	// create ticker when fisrt tick called
	// Init first beat
	b.addOneTick(time.Now())
	b.status = BeatRun
	go b.listen()
	go b.tick()
}

// Stop Export, blocking
func (b *Beat) Stop() {
	b.status = BeatStop

	// call context cancelfunc to notify tick and listen
	b.cancelFunc()

	// 等待Run退出
	// fmt.Println("wait")
	b.wg.Wait()
	// fmt.Println("wait over")

	// close chan
	close(b.tickChan)
	// fmt.Println("chan closed")
	b.ticker.Stop()
}
