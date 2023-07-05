package hiface

import "time"

type IHeartBeatChecker interface {
	ResetHeartBreaker()
	StopHeartBreaker()
	TickHeartBreaker() <-chan time.Time
}
