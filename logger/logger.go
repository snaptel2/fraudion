package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"io/ioutil"
)

const (
	constSignature = "FRAUDION"
	constFlag      = log.Ldate | log.Ltime | log.Lshortfile

	// ConstLoggerLevelTrace ...
	ConstLoggerLevelTrace = "trace"
	// ConstLoggerLevelInfo ...
	ConstLoggerLevelInfo = "info"
	// ConstLoggerLevelWarning ...
	ConstLoggerLevelWarning = "warning"
	// ConstLoggerLevelError ...
	ConstLoggerLevelError = "error"
)

// Log ...
var Log *Logger

// Logger ...
type Logger struct {
	Level   string
	Handles loggerHandles
}

type loggerHandles struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func init() {
	Log = new(Logger)
	Log.Level = ConstLoggerLevelInfo
	Log.SetHandles(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

// SetHandles ...
func (logger *Logger) SetHandles(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	logger.Handles.Trace = log.New(traceHandle, fmt.Sprintf("%s TRACE: ", constSignature), constFlag)
	logger.Handles.Info = log.New(infoHandle, fmt.Sprintf("%s INFO: ", constSignature), constFlag)
	logger.Handles.Warning = log.New(warningHandle, fmt.Sprintf("%s WARNING: ", constSignature), constFlag)
	logger.Handles.Error = log.New(errorHandle, fmt.Sprintf("%s ERROR: ", constSignature), constFlag)
}

// Log ...
func (logger *Logger) Write(level string, message string, fatal bool) {
	level = strings.ToLower(level)
	switch level {
	case ConstLoggerLevelTrace:
		logger.Handles.Trace.Println(message)
	case ConstLoggerLevelInfo:
		logger.Handles.Info.Println(message)
	case ConstLoggerLevelWarning:
		logger.Handles.Warning.Println(message)
	case ConstLoggerLevelError:
		if fatal {
			logger.Handles.Error.Fatalln(message)
		} else {
			logger.Handles.Error.Println(message)
		}
	default:
		logger.Handles.Info.Println(message)
	}
}
