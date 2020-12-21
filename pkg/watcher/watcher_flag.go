package watcher

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	DefaultSecond = 2
)

type WatcherFlags struct {
	Second int
	Watch  bool
}

func (w *WatcherFlags) AddFlags(cmd *cobra.Command, usage string) {
	cmd.Flags().BoolVarP(&w.Watch, "watch", "", false, usage)
	cmd.Flags().IntVarP(&w.Second, "watch-second", "", DefaultSecond, "Refresh seconds")
}

func (w *WatcherFlags) ToWatch(fn WatcherFunc) {
	if w.Watch {
		if w.Second <= 0 {
			w.Second = DefaultSecond
		}
		watch := &Watcher{Second: time.Duration(w.Second)}
		watch.Run(fn)
	}
}

func NewWatcherFlags() *WatcherFlags {
	return &WatcherFlags{}
}
