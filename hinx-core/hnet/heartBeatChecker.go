package hnet

import (
	"hinx/hinx-core/hconf"
	"time"
)

type heartBeatChecker struct {
	*time.Timer
}

func (h *heartBeatChecker) ResetHeartBreaker() {
	h.Timer.Reset(time.Duration(hconf.GlobalObject.Heartbeat) * time.Millisecond)
}

func (h *heartBeatChecker) StopHeartBreaker() {
	h.Timer.Stop()
}

func (h *heartBeatChecker) TickHeartBreaker() <-chan time.Time {
	return h.Timer.C
}
