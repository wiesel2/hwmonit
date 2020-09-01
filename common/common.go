package common

import (
	"os"
	"path"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger Export
var logger *zap.SugaredLogger

func initlogger() {
	cfg := GetConfig()
	hook := lumberjack.Logger{
		Filename:   path.Join(cfg["log-dir"], cfg["log"]), // 日志文件路径
		MaxSize:    30,                                    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 20,                                    // 日志文件最多保存多少个备份
		MaxAge:     7,                                     // 文件最多保存多少天
		Compress:   true,                                  // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	// filed := zap.Fields(zap.String("serviceName", "codesync"))
	// 构造日志
	l := zap.New(core, caller, development)
	logger = l.Sugar()
}

// GetLogger Export
func GetLogger() *zap.SugaredLogger {
	return logger
}

var once sync.Once

func init() {
	once.Do(func() {
		readArgs()
		readConfig()
		initlogger()
	})
}
