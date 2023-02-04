package gonuts

import (
	"time"
)

func TimeFromUnixTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// this is just to remember that javascript Date.now() converts like this
func TimeFromJSTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, 0)
}

// this is just to remember that javascript Date.now() converts like this
func TimeToJSTimestamp(t time.Time) int64 {
	return t.UnixMilli()
}
