package gonuts

import "time"

/*
	intervalChannel := Interval(time.Duration(time.Second*1), func() { nuts.L.Debugf("tick ", time.Now()) }, true)
	intervalChannel.Stop()
*/

func Interval(theCall func(), dur time.Duration, runFuncImmediately bool) *time.Ticker {
	ticker := time.NewTicker(dur)
	go func() {
		for range ticker.C {
			theCall()
		}
	}()
	if runFuncImmediately {
		theCall()
	}
	return ticker
}
