package logger

import "sync"

var (
	_log ILogger
	once sync.Once
)

// Init 初始化日志
func Init(cfg *Config) (ILogger, error) {
	var err error
	once.Do(func() {
		_log, err = NewLogger(cfg)
	})
	return _log, err
}

// Debug 记录debug级别日志
func Debug(args ...any) {
	_log.Debug(args...)
}

// Debugf 记录debug级别日志（格式化）
func Debugf(format string, args ...any) {
	_log.Debugf(format, args...)
}

// Info 记录info级别日志
func Info(args ...any) {
	_log.Info(args...)
}

// Infof 记录info级别日志（格式化）
func Infof(format string, args ...any) {
	_log.Infof(format, args...)
}

// Warn 记录warn级别日志
func Warn(args ...any) {
	_log.Warn(args...)
}

// Warnf 记录warn级别日志（格式化）
func Warnf(format string, args ...any) {
	_log.Warnf(format, args...)
}

// Error 记录error级别日志
func Error(args ...any) {
	_log.Error(args...)
}

// Errorf 记录error级别日志（格式化）
func Errorf(format string, args ...any) {
	_log.Errorf(format, args...)
}

// Fatal 记录fatal级别日志
func Fatal(args ...any) {
	_log.Fatal(args...)
}

// Fatalf 记录fatal级别日志（格式化）
func Fatalf(format string, args ...any) {
	_log.Fatalf(format, args...)
}

// Panic 记录panic级别日志
func Panic(args ...any) {
	_log.Panic(args...)
}

// Panicf 记录panic级别日志（格式化）
func Panicf(format string, args ...any) {
	_log.Panicf(format, args...)
}
