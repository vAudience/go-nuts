package gonuts

import (
	"time"
)

func TimeFromUnixTimestamp(timestamp int64) time.Time {
	tm := time.Unix(timestamp, 0)
	return tm
}
