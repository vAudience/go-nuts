package gonuts

import "go.uber.org/zap/zapcore"

var L = Init_Logger(zapcore.DebugLevel, "unknown", false, "logs/")
