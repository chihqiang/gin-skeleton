package logger

// 导入必要的包
import (
	"go.uber.org/zap"         // Zap 日志库
	"go.uber.org/zap/zapcore" // Zap 核心组件
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
)

// 日志级别常量定义
const (
	DebugLevel = "debug" // 调试级别
	InfoLevel  = "info"  // 信息级别
	WarnLevel  = "warn"  // 警告级别
	ErrorLevel = "error" // 错误级别
	FatalLevel = "fatal" // 致命错误级别
	PanicLevel = "panic" // 恐慌级别
)

// ILogger 日志记录器接口
// 定义了各种日志级别的记录方法
type ILogger interface {
	// Debug 记录 debug 级别日志
	Debug(args ...any)
	// Debugf 记录带格式的 debug 级别日志
	Debugf(format string, args ...any)
	// Info 记录 info 级别日志
	Info(args ...any)
	// Infof 记录带格式的 info 级别日志
	Infof(format string, args ...any)
	// Warn 记录 warn 级别日志
	Warn(args ...any)
	// Warnf 记录带格式的 warn 级别日志
	Warnf(format string, args ...any)
	// Error 记录 error 级别日志
	Error(args ...any)
	// Errorf 记录带格式的 error 级别日志
	Errorf(format string, args ...any)
	// Fatal 记录 fatal 级别日志
	Fatal(args ...any)
	// Fatalf 记录带格式的 fatal 级别日志
	Fatalf(format string, args ...any)
	// Panic 记录 panic 级别日志
	Panic(args ...any)
	// Panicf 记录带格式的 panic 级别日志
	Panicf(format string, args ...any)
}

// Config 日志配置结构体
// 包含日志系统的各种配置参数
type Config struct {
	Level      string `yaml:"level" json:"level,omitempty"`             // 日志级别
	Path       string `yaml:"path" json:"path,omitempty"`               // 日志文件路径,没有定义就走控制台
	MaxSize    int    `yaml:"max_size" json:"max_size,omitempty"`       // 单个日志文件最大大小（MB）
	MaxBackups int    `yaml:"max_backups" json:"max_backups,omitempty"` // 保留的最大备份文件数
	MaxAge     int    `yaml:"max_age" json:"max_age,omitempty"`         // 日志文件最大保存天数
	Compress   bool   `yaml:"compress" json:"compress,omitempty"`       // 是否压缩备份文件
	Format     string `yaml:"format" json:"format,omitempty"`           // 日志格式（json 或 text）
}

// Logger 日志记录器实现
// 包装了 zap.SugaredLogger 提供日志记录功能
type Logger struct {
	sync.RWMutex
	log *zap.SugaredLogger
}

// NewLogger 创建 Logger
func NewLogger(cfg *Config) (*Logger, error) {
	var cores []zapcore.Core
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 时间字段名
		LevelKey:       "level",                        // 级别字段名
		NameKey:        "logger",                       // 日志器名称字段名
		CallerKey:      "caller",                       // 调用者字段名
		MessageKey:     "msg",                          // 消息字段名
		StacktraceKey:  "stacktrace",                   // 堆栈跟踪字段名
		LineEnding:     zapcore.DefaultLineEnding,      // 行结束符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 级别编码器（小写）
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 时间编码器（ISO8601格式）
		EncodeDuration: zapcore.SecondsDurationEncoder, // 持续时间编码器
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 调用者编码器（短格式）
	}
	var encoder zapcore.Encoder
	if strings.ToLower(cfg.Format) == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	// 日志级别
	level := zapcore.InfoLevel
	switch strings.ToLower(cfg.Level) {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	case FatalLevel:
		level = zapcore.FatalLevel
	case PanicLevel:
		level = zapcore.PanicLevel
	}
	if cfg.Path != "" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   cfg.Path,       // 日志文件路径（完整文件名，如 runtime/log/app.log）
			MaxSize:    cfg.MaxSize,    // 单个日志文件最大大小（单位：MB），超过会切分新文件
			MaxBackups: cfg.MaxBackups, // 保留的最大历史备份文件数量，超过会删除最旧的
			MaxAge:     cfg.MaxAge,     // 日志文件最大保存天数，超过天数会删除旧文件
			Compress:   cfg.Compress,   // 是否压缩备份文件，启用后旧日志会生成 .gz 文件
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), level))
	} else {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}
	// 合并核心
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return &Logger{log: logger}, nil
}

// Debug 记录debug级别日志
func (l *Logger) Debug(args ...any) {
	l.log.Debug(args...)
}

// Debugf 记录debug级别日志（格式化）
func (l *Logger) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

// Info 记录info级别日志
func (l *Logger) Info(args ...any) {
	l.log.Info(args...)
}

// Infof 记录info级别日志（格式化）
func (l *Logger) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

// Warn 记录warn级别日志
func (l *Logger) Warn(args ...any) {
	l.log.Warn(args...)
}

// Warnf 记录warn级别日志（格式化）
func (l *Logger) Warnf(format string, args ...any) {
	l.log.Warnf(format, args...)
}

// Error 记录error级别日志
func (l *Logger) Error(args ...any) {
	l.log.Error(args...)
}

// Errorf 记录error级别日志（格式化）
func (l *Logger) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

// Fatal 记录fatal级别日志
func (l *Logger) Fatal(args ...any) {
	l.log.Fatal(args...)
}

// Fatalf 记录fatal级别日志（格式化）
func (l *Logger) Fatalf(format string, args ...any) {
	l.log.Fatalf(format, args...)
}

// Panic 记录panic级别日志
func (l *Logger) Panic(args ...any) {
	l.log.Panic(args...)
}

// Panicf 记录panic级别日志（格式化）
func (l *Logger) Panicf(format string, args ...any) {
	l.log.Panicf(format, args...)
}
