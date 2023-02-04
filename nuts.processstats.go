package gonuts

import (
	"fmt"
	"os"
	"runtime"
)

func PrintMemoryUsage() bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	content := fmt.Sprintf("[MemoryUsage] PID=(%d) | Alloc = %v MiB | TotalAlloc = %v MiB | Sys = %v MiB | NumGC = %v", os.Getpid(), bToMb(int64(m.Alloc)), bToMb(int64(m.TotalAlloc)), bToMb(int64(m.Sys)), m.NumGC)
	L.Debugf(content)
	return true
}

func bToMb(b int64) int64 {
	return b / 1024 / 1024
}
