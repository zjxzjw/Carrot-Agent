package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	mu         sync.Mutex
	level      Level
	output     *os.File
	timeFormat string
}

var (
	defaultLogger *Logger
	once         sync.Once
)

func init() {
	once.Do(func() {
		defaultLogger = &Logger{
			level:      INFO,
			output:     os.Stdout,
			timeFormat: "2006-01-02 15:04:05",
		}
	})
}

func NewLogger(level Level, output *os.File) *Logger {
	return &Logger{
		level:      level,
		output:     output,
		timeFormat: "2006-01-02 15:04:05",
	}
}

func Default() *Logger {
	return defaultLogger
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetOutput(output *os.File) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format(l.timeFormat)
	_, file, line, ok := runtime.Caller(2)
	var location string
	if ok {
		location = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	} else {
		location = "unknown"
	}

	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "[%s] [%s] [%s] %s\n", timestamp, level.String(), location, message)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

func Debug(format string, args ...interface{}) {
	defaultLogger.log(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.log(INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.log(WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.log(ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.log(FATAL, format, args...)
	os.Exit(1)
}

func ParseLevel(levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}