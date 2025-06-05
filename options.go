/*
 * @Author: yeying
 * @Date: 2025-06-03 21:53:59
 * @LastEditTime: 2025-06-03 22:49:53
 * @FilePath: \dark\options.go
 * @Description:
 */
package dark

type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

type Options struct {
	Addr        string
	IsWebsocket bool
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
