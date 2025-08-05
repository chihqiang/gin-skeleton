package logger

// 导入必要的包
import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"         // Zap 日志库
	"go.uber.org/zap/zapcore" // Zap 核心组件
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
	Level      string `yaml:"level"`       // 日志级别
	Path       string `yaml:"path"`        // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小（MB）
	MaxBackups int    `yaml:"max_backups"` // 保留的最大备份文件数
	MaxAge     int    `yaml:"max_age"`     // 日志文件最大保存天数
	Compress   bool   `yaml:"compress"`    // 是否压缩备份文件
	Format     string `yaml:"format"`      // 日志格式（json 或 console）
}

// Logger 日志记录器实现
// 包装了 zap.SugaredLogger 提供日志记录功能
type Logger struct {
	log *zap.SugaredLogger // Zap 日志记录器实例
}

// NewLogger 创建一个新的日志记录器实例
// 根据配置初始化日志系统
// 参数 cfg: 日志配置
// 返回值: 日志记录器实例或错误
func NewLogger(cfg *Config) (*Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return nil, fmt.Errorf("create log directory failed: %w", err)
	}

	// 解析日志级别
	var level zapcore.Level
	switch cfg.Level {
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
	default:
		level = zapcore.InfoLevel // 默认使用 info 级别
	}

	// 编码配置
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

	// 根据配置选择编码器（JSON 或 Console）
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 按天创建日志文件：logs/2025-08-06.log
	filename := filepath.Join(cfg.Path, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file failed: %w", err)
	}

	// 创建写入器
	writer := zapcore.AddSync(logFile)

	// 创建日志核心
	// 核心组件负责将日志条目编码并写入到指定的写入器
	core := zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(level))

	// 创建日志器
	// zap.AddCaller() 添加调用者信息
	// zap.AddStacktrace(zapcore.ErrorLevel) 为错误级别及以上添加堆栈跟踪
	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// 返回日志记录器实例
	return &Logger{
		log: zapLogger.Sugar(),
	}, nil
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
