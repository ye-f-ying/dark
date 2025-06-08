package logger

import "github.com/rs/zerolog"

// Level 定义日志级别类型
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
	//LevelPanic Level = "panic"
	LevelTrace Level = "trace"
)

// 级别映射
var levelMap = map[Level]zerolog.Level{
	LevelDebug: zerolog.DebugLevel,
	LevelInfo:  zerolog.InfoLevel,
	LevelWarn:  zerolog.WarnLevel,
	LevelError: zerolog.ErrorLevel,
	LevelFatal: zerolog.FatalLevel,
	//LevelPanic: zerolog.PanicLevel,
	LevelTrace: zerolog.TraceLevel,
}

// levelWriter 实现按级别过滤的写入器
type levelWriter struct {
	zerolog.LevelWriter
	level   zerolog.Level
	include bool // 是否包含该级别（true: >=level, false: ==level）
}

// WriteLevel 实现 zerolog.LevelWriter 接口
func (w *levelWriter) WriteLevel(lvl zerolog.Level, p []byte) (int, error) {
	if w.include {
		if lvl >= w.level { // 包含该级别及以上
			return w.LevelWriter.Write(p)
		}
	} else {
		if lvl == w.level { // 仅包含该级别
			return w.LevelWriter.Write(p)
		}
	}
	return len(p), nil // 忽略不符合条件的级别
}
