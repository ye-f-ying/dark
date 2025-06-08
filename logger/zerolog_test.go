/*
 * @Author: yeying
 * @Date: 2025-06-08 18:08:49
 * @LastEditTime: 2025-06-08 18:19:22
 * @FilePath: \dark\logger\zerolog_test.go
 * @Description:
 */
package logger

import (
	"testing"
)

func TestConsle(t *testing.T) {
	// 创建默认配置的 logger（输出到控制台，人类可读格式）
	log, err := NewLogger(nil)
	if err != nil {
		panic(err)
	}

	// 基本使用
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")

	// 带格式化的日志
	userID := 123
	log.Infof("Infof用户 %d 登录成功", userID)
}

func TestFile(t *testing.T) {
	// 配置：输出到文件，禁用控制台输出
	cfg := &LoggerConfig{
		Console: false,    // 禁用控制台
		LogDir:  "./logs", // 日志文件目录
		//Level:   "info",   // 日志级别
	}

	log, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}

	// 日志将写入 ./logs/app-YYYY-MM-DD.log
	log.Info("Info")
	log.Error("Error")
}

func TestFileConsle(t *testing.T) {
	// 配置：同时输出到控制台和文件
	cfg := &LoggerConfig{
		Console:      true,     // 启用控制台
		ConsoleHuman: true,     // 控制台使用人类可读格式
		LogDir:       "./logs", // 日志文件目录
		Level:        "debug",  // 日志级别
	}

	log, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}

	// 控制台和文件都会收到这些日志
	log.Debug("Debug")
	log.Info("Info")
}

func TestTypeFileConsle(t *testing.T) {
	// 配置：同时输出到控制台和文件
	cfg := &LoggerConfig{
		Console:      true,     // 启用控制台
		ConsoleHuman: true,     // 控制台使用人类可读格式
		LogDir:       "./logs", // 日志文件目录
		Level:        "trace",  // 日志级别
		SplitByLevel: false,    //是否按级别拆分 true 是
	}

	log, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}

	// 控制台和文件都会收到这些日志
	log.Debug("Debug")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	log.Trace("Trace")
}
