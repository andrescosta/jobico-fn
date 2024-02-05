package executor

import "time"

type ticker interface {
	Chan() <-chan time.Time
	Stop()
	Tick()
}

type timeBasedTicker struct {
	ticker *time.Ticker
}

func (t *timeBasedTicker) Chan() <-chan time.Time {
	return t.ticker.C
}

func (t *timeBasedTicker) Stop() {
	t.ticker.Stop()
}

func (t *timeBasedTicker) Tick() {
}

type channelBasedTicker struct {
	c chan time.Time
}

func (t *channelBasedTicker) Chan() <-chan time.Time {
	return t.c
}

func (t *channelBasedTicker) Stop() {
	close(t.c)
}

func (t *channelBasedTicker) Tick() {
	t.c <- time.Now()
}
