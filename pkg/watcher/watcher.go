package watcher

import (
	"time"
)

type Watcher struct {
	Second time.Duration
}

type WatcherFunc func(chan bool)

func (w *Watcher) Run(fn WatcherFunc) {
	d := time.Duration(w.Second * time.Second)
	done := make(chan bool, 1)
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-done:
			return
		case <-t.C:
			fn(done)
		}
	}
}
