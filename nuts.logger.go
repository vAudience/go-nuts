package gonuts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	GO_NUTS_LOGGER_CONFIG      = "GO_NUTS_LOGGER_CONFIG"
	GO_NUTS_LOGGER_CONFIG_PROD = "prod"
)

// CHECK https://stackoverflow.com/questions/68472667/how-to-log-to-stdout-or-stderr-based-on-log-level-using-uber-go-zap
func Init_Logger(targetLevel zapcore.Level, instanceId string, log2file bool, logfilePath string) *zap.SugaredLogger {
	var LogConfig zap.Config

	// Set default log config
	LogConfig = zap.NewDevelopmentConfig()
	LogConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	LogConfig.EncoderConfig.EncodeTime = SyslogTimeEncoder
	LogConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	LogConfig.Level = zap.NewAtomicLevelAt(targetLevel)

	// Set production log config based on environment variable
	loggerConfig, ok := os.LookupEnv(GO_NUTS_LOGGER_CONFIG)
	if ok && loggerConfig == GO_NUTS_LOGGER_CONFIG_PROD {
		LogConfig = zap.NewProductionConfig()
		LogConfig.Level = zap.NewAtomicLevelAt(targetLevel)
	}

	if log2file && logfilePath != "" {
		logfileName := logfilePath + "log_" + time.Now().Format("2006-01-02T15:04:05Z07:00") + "_" + instanceId + ".txt"
		LogConfig.OutputPaths = append(LogConfig.OutputPaths, logfileName)
		fmt.Printf("[nuts.logger] adding logfile: (%s)", logfileName)
	}

	logger, err := LogConfig.Build()
	if err != nil {
		fmt.Printf("[nuts.logger] ERROR! failed to create logger PANIC! \n%s", err)
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

func SetLoglevel(loglevel string, instanceId string, log2file bool, logfilePath string) {
	switch loglevel {
	case "DEBUG":
		L = Init_Logger(zap.DebugLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to DEBUG.")
	case "INFO":
		L = Init_Logger(zap.InfoLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to INFO.")
	case "WARN":
		L = Init_Logger(zap.WarnLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to WARN.")
	case "ERROR":
		L = Init_Logger(zap.ErrorLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to ERROR.")
	case "FATAL":
		L = Init_Logger(zap.FatalLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to FATAL.")
	case "PANIC":
		L = Init_Logger(zap.PanicLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to PANIC.")
	default:
		L = Init_Logger(zap.DebugLevel, instanceId, log2file, logfilePath)
		fmt.Println("[SetLoglevel] LogLevel set to DEFAULT (DEBUG).")
	}
}

func GetPrettyJson(object any) (pretty string) {
	pretty = ""
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return "failed to marshal to json :("
	}
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, jsonBytes, "", "\t")
	pretty = prettyJSON.String()
	return pretty
}
