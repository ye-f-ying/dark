/*
 * @Author: yeying
 * @Date: 2025-06-03 21:53:59
 * @LastEditTime: 2025-06-10 21:58:59
 * @FilePath: \dark\options.go
 * @Description:
 */
package dark

import (
	"dark/logger"

	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
)

type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

type Options struct {
	Addr            string
	IsWebsocket     bool
	Upgrader        websocket.Upgrader
	isUpgrader      bool
	LogConsole      *bool         // 是否输出到控制台
	LogConsoleHuman *bool         // 控制台是否为人类可读格式
	LogDir          string        // 日志文件目录，默认 "" 表示不写入文件
	LogSplitByLevel *bool         // 是否按级别分割日志
	LogLevel        string        // 日志级别
	GoPool          IGoPool       //协程池
	isAnts          *bool         //是否使用ants协程池
	antsOpt         []ants.Option //
	antsSize        int
}

/**
 * @description: 设置地址
 * @param {string} addr
 * @return {*}
 */
func WithAddr(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

/**
 * @description:是否启用websocket
 * @param {bool} IsWebsocket
 * @return {*}
 */
func WithIsWebsocket(IsWebsocket bool) Option {
	return func(opts *Options) {
		opts.IsWebsocket = IsWebsocket
	}
}

/**
 * @description: 设置upgrader
 * @param {websocket.Upgrader} upgrader
 * @return {*}
 */
func WithUpgrader(upgrader websocket.Upgrader) Option {
	return func(opts *Options) {
		opts.Upgrader = upgrader
		opts.isUpgrader = true
	}
}

/**
 * @description: 是否输出到控制台
 * @param {bool} console
 * @return {*}
 */
func WithLogConsole(console bool) Option {
	return func(opts *Options) {
		opts.LogConsole = &console
	}
}

/**
 * @description: 控制台是否为人类可读格式
 * @param {bool} consoleHuman
 * @return {*}
 */
func WithLogConsoleHuman(consoleHuman bool) Option {
	return func(opts *Options) {
		opts.LogConsoleHuman = &consoleHuman
	}
}

/**
 * @description: 日志文件目录，默认 "" 表示不写入文件
 * @param {string} dir
 * @return {*}
 */
func WithLogDir(dir string) Option {
	return func(opts *Options) {
		opts.LogDir = dir
	}
}

/**
 * @description: 是否按级别分割日志
 * @param {bool} spl
 * @return {*}
 */
func WithLogSplitByLevel(spl bool) Option {
	return func(opts *Options) {
		opts.LogSplitByLevel = &spl
	}
}

/**
 * @description: 日志级别
 * @param {logger.Level} level
 * @return {*}
 */
func WithLogLevel(level logger.Level) Option {
	return func(opts *Options) {
		opts.LogLevel = string(level)
	}
}

/**
 * @description: 自定义使用协程池对象
 * @param {IGoPool} p
 * @return {*}
 */
func WithGoPool(p IGoPool) Option {
	return func(opts *Options) {
		opts.GoPool = p
	}
}

/**
 * @description: 是否使用ant协程池
 * @param {bool} p
 * @return {*}
 */
func WithAnts(size int, antOpts ...ants.Option) Option {
	return func(opts *Options) {
		isAnts := true
		opts.isAnts = &isAnts
		opts.antsOpt = antOpts
		opts.antsSize = size
	}
}
