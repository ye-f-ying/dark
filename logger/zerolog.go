package logger

import (
	"context"
	"dark/logger/dlog"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/zerolog"
)

// LoggerConfig 控制台输出配置
type LoggerConfig struct {
	Console      bool   // 是否输出到控制台
	ConsoleHuman bool   // 控制台是否为人类可读格式
	LogDir       string // 日志文件目录，默认 "" 表示不写入文件
	SplitByLevel bool   // 是否按级别分割日志
	Level        string // 日志级别
}

func GetDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Console:      true,
		ConsoleHuman: true,
		SplitByLevel: false,
		Level:        string(LevelTrace),
	}

}

// 创建按级别分割的写入器
func createLevelWriters(logDir string) ([]io.Writer, error) {
	var writers []io.Writer

	for level, zl := range levelMap {
		// 构建级别对应的日志文件名（如 error.log）
		logType := string(level)
		path := fmt.Sprintf("%s/%s-%%Y-%%m-%%d.log", logDir, logType)

		// 创建日志文件写入器（禁用符号链接）
		writer, err := rotatelogs.New(
			path,
			rotatelogs.WithMaxAge(7*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			return nil, fmt.Errorf("创建 %s 写入器失败: %w", logType, err)
		}

		// 创建级别过滤写入器（默认仅写入当前级别）
		levelWriter := &levelWriter{
			LevelWriter: zerolog.LevelWriterAdapter{Writer: writer},
			level:       zl,
			include:     false, // 仅写入当前级别（如需包含以上级别，设为 true）
		}

		writers = append(writers, levelWriter)
	}

	return writers, nil
}

// ZeroLogger 实现 dlog.FullLogger 接口
type ZeroLogger struct {
	logger        zerolog.Logger
	level         zerolog.Level
	logToFile     bool
	consoleWriter io.Writer
	fileWriter    io.Writer
	mu            sync.RWMutex
}

// newRotatingWriter 创建每天分割的 rotatelogs writer
func newRotatingWriter(logDir, logType string) (io.Writer, error) {
	// 确保目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	path := fmt.Sprintf("%s/%s-%%Y-%%m-%%d.log", logDir, logType)
	//link := fmt.Sprintf("%s/%s.log", logDir, logType)

	writer, err := rotatelogs.New(
		path,
		//rotatelogs.WithLinkName(link),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("创建日志写入器失败: %w", err)
	}

	return writer, nil
}

/**
 * @description: NewLogger 创建日志实例并返回 dlog.FullLogger 接口
 * @param {*LoggerConfig} cfg
 * @return {*}
 */
func NewLogger(cfg *LoggerConfig) (dlog.FullLogger, error) {
	if cfg == nil {
		cfg = GetDefaultLoggerConfig()
	}
	console := cfg.Console
	consoleHuman := cfg.ConsoleHuman
	logDir := cfg.LogDir
	splitByLevel := cfg.SplitByLevel
	levelStr := cfg.Level

	// 解析日志级别
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.DebugLevel // 默认级别
	}

	var writers []io.Writer

	if splitByLevel && logDir != "" {
		// 按级别分割日志
		levelWriters, err := createLevelWriters(logDir)
		if err != nil {
			return nil, err
		}
		writers = levelWriters
	} else if logDir != "" {
		// 非分割模式：创建综合日志文件
		fileWriter, err := newRotatingWriter(logDir, "dark_app")
		if err != nil {
			return nil, err
		}
		writers = append(writers, fileWriter)
	}

	// 添加控制台写入器
	if console {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    !consoleHuman,
		}
		writers = append(writers, consoleWriter)
	}

	// 处理无输出的情况
	if len(writers) <= 0 {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    !consoleHuman,
		})
	}

	// 构建 logger
	logger := zerolog.New(zerolog.MultiLevelWriter(writers...)).
		Level(level).
		With().Timestamp().Logger()

	return &ZeroLogger{
		logger: logger,
		level:  level,
		mu:     sync.RWMutex{},
	}, nil
}

func (z *ZeroLogger) Debug(args ...interface{}) {
	z.logger.Debug().Msg(fmt.Sprint(args...))
}

func (z *ZeroLogger) Error(args ...interface{}) {
	z.logger.Error().Msg(fmt.Sprint(args...))
}

func (z *ZeroLogger) Fatal(args ...interface{}) {
	z.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (z *ZeroLogger) Info(args ...interface{}) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

func (z *ZeroLogger) Notice(args ...interface{}) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

func (z *ZeroLogger) Trace(args ...interface{}) {
	z.logger.Trace().Msg(fmt.Sprint(args...))
}

func (z *ZeroLogger) Warn(args ...interface{}) {
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

/**
 * @description: 实现 dlog.Logger 接口的 SetOutput 方法
 * @param {io.Writer} w
 * @return {*}
 */
func (z *ZeroLogger) SetOutput(w io.Writer) {
	z.mu.Lock()
	defer z.mu.Unlock()

	// 更新控制台写入器
	z.consoleWriter = w

	// 重新构建写入器列表
	var writers []io.Writer
	if z.logToFile && z.fileWriter != nil {
		writers = append(writers, z.fileWriter)
	}
	writers = append(writers, z.consoleWriter)

	// 创建新的 logger 并设置输出
	z.logger = z.logger.Hook(nil).Output(zerolog.MultiLevelWriter(writers...)).Level(z.level)
}

func (z *ZeroLogger) Debugf(format string, args ...interface{}) {
	z.logger.Debug().Msgf(format, args...)
}

func (z *ZeroLogger) Infof(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...)
}

func (z *ZeroLogger) Noticef(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...) // 使用 Info 替代 Notice
}

func (z *ZeroLogger) Warnf(format string, args ...interface{}) {
	z.logger.Warn().Msgf(format, args...)
}

func (z *ZeroLogger) Errorf(format string, args ...interface{}) {
	z.logger.Error().Msgf(format, args...)
}

func (z *ZeroLogger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Msgf(format, args...)
}

func (z *ZeroLogger) Tracef(format string, args ...interface{}) {
	z.logger.Trace().Msgf(format, args...)
}

func (z *ZeroLogger) CtxDebugf(ctx context.Context, format string, args ...interface{}) {
	z.withContext(ctx).Debug().Msgf(format, args...)
}

func (z *ZeroLogger) CtxInfof(ctx context.Context, format string, args ...interface{}) {
	z.withContext(ctx).Info().Msgf(format, args...)
}

func (z *ZeroLogger) CtxNoticef(ctx context.Context, format string, args ...interface{}) {
	z.withContext(ctx).Info().Msgf(format, args...) // 使用 Info 替代 Notice
}

func (z *ZeroLogger) CtxWarnf(ctx context.Context, format string, args ...interface{}) {
	z.withContext(ctx).Warn().Msgf(format, args...)
}

func (z *ZeroLogger) CtxErrorf(ctx context.Context, format string, args ...interface{}) {
	z.withContext(ctx).Error().Msgf(format, args...)
}

func (z *ZeroLogger) CtxFatalf(ctx context.Context, format string, args ...interface{}) {
	z.withContext(ctx).Fatal().Msgf(format, args...)
}

func (z *ZeroLogger) CtxTracef(ctx context.Context, format string, args ...interface{}) {
	z.withContext(ctx).Trace().Msgf(format, args...)
}

/**
 * @description: 辅助方法：从上下文中提取 request_id
 * @param {context.Context} ctx
 * @return {*}
 */
func (z *ZeroLogger) withContext(ctx context.Context) *zerolog.Logger {
	if ctx == nil {
		return &z.logger
	}

	if reqID, ok := ctx.Value("X-Request-ID").(string); ok {
		logger := z.logger.With().Str("request_id", reqID).Logger()
		return &logger
	}
	return &z.logger
}

/**
 * @description:设置日志等级
 * @param {dlog.Level} level
 * @return {*}
 */
func (z *ZeroLogger) SetLevel(level dlog.Level) {
	z.mu.Lock()
	defer z.mu.Unlock()

	switch level {
	case dlog.LevelDebug:
		z.level = zerolog.DebugLevel
	case dlog.LevelInfo:
		z.level = zerolog.InfoLevel
	case dlog.LevelWarn:
		z.level = zerolog.WarnLevel
	case dlog.LevelError:
		z.level = zerolog.ErrorLevel
	case dlog.LevelFatal:
		z.level = zerolog.FatalLevel
	default:
		z.level = zerolog.InfoLevel
	}
	z.logger = z.logger.Level(z.level)
}

func (z *ZeroLogger) Level() dlog.Level {
	z.mu.RLock()
	defer z.mu.RUnlock()

	switch z.level {
	case zerolog.DebugLevel:
		return dlog.LevelDebug
	case zerolog.InfoLevel:
		return dlog.LevelInfo
	case zerolog.WarnLevel:
		return dlog.LevelWarn
	case zerolog.ErrorLevel:
		return dlog.LevelError
	case zerolog.FatalLevel:
		return dlog.LevelFatal
	default:
		return dlog.LevelInfo
	}
}
