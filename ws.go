package dark

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
)

type WEBSOCKET_CONTEXT_TYPE string

const WEBSOCKET_CONTEXT_TYPE_KEY WEBSOCKET_CONTEXT_TYPE = "_WebSocketMessageTypeKey_"

// WebSocket 帧类型
const (
	FrameContinuation = 0x0
	FrameText         = 0x1
	FrameBinary       = 0x2
	FrameClose        = 0x8
	FramePing         = 0x9
	FramePong         = 0xA
)

// WebSocketFrameHeader 表示WebSocket帧头
type WebSocketFrameHeader struct {
	FIN       bool
	RSV1      bool
	RSV2      bool
	RSV3      bool
	OpCode    byte
	Masked    bool
	Length    uint64
	MaskKey   [4]byte
	HeaderLen int // 帧头长度（字节）
}

// ParseWebSocketFrame 解析WebSocket帧
// 返回：帧头, 帧数据, 剩余数据, 错误
/**
 * @description: 解析WebSocket帧
 * @param {[]byte} data
 * @return {*}
 */
func ParseWebSocketFrame(data []byte) (int, *WebSocketFrameHeader, []byte, error) {
	// 检查是否有足够的数据读取基本头部（至少2字节）
	if len(data) < 2 {
		return 0, nil, nil, nil
	}

	// 解析第一个字节
	header := &WebSocketFrameHeader{
		FIN:    (data[0] & 0x80) != 0,
		RSV1:   (data[0] & 0x40) != 0,
		RSV2:   (data[0] & 0x20) != 0,
		RSV3:   (data[0] & 0x10) != 0,
		OpCode: data[0] & 0x0F,
	}

	// 解析第二个字节
	header.Masked = (data[1] & 0x80) != 0
	payloadLen := data[1] & 0x7F
	header.HeaderLen = 2 // 基本头部大小

	// 解析扩展长度
	switch payloadLen {
	case 126:
		// 16位扩展长度
		if len(data) < header.HeaderLen+2 {
			return 0, nil, nil, nil
		}
		header.Length = uint64(binary.BigEndian.Uint16(data[header.HeaderLen : header.HeaderLen+2]))
		header.HeaderLen += 2
	case 127:
		// 64位扩展长度
		if len(data) < header.HeaderLen+8 {
			return 0, nil, nil, nil
		}
		header.Length = binary.BigEndian.Uint64(data[header.HeaderLen : header.HeaderLen+8])
		// 检查长度是否过大（避免内存溢出）
		if header.Length > 1024*1024*10 { // 限制最大10MB
			return 0, nil, nil, fmt.Errorf("消息长度过大: %d", header.Length)
		}
		header.HeaderLen += 8
	default:
		// 直接使用7位长度
		header.Length = uint64(payloadLen)
	}

	// 解析掩码密钥
	if header.Masked {
		if len(data) < header.HeaderLen+4 {
			return 0, nil, nil, nil
		}
		copy(header.MaskKey[:], data[header.HeaderLen:header.HeaderLen+4])
		header.HeaderLen += 4
	}

	// 检查是否有足够的数据读取完整的帧
	totalFrameLen := header.HeaderLen + int(header.Length)
	if len(data) < totalFrameLen {
		return 0, nil, nil, nil
	}

	// 提取帧数据
	frameData := data[header.HeaderLen:totalFrameLen]

	// 应用掩码（如果有）
	if header.Masked {
		for i := range frameData {
			frameData[i] ^= header.MaskKey[i%4]
		}
	}
	// 提取剩余数据
	//remainingData := data[totalFrameLen:]
	return totalFrameLen, header, frameData, nil
}

/**
 * @description: 发布websocket帧数据
 * @param {byte} opcode
 * @param {[]byte} data
 * @param {bool} fin
 * @param {bool} mask
 * @return {*}
 */
func BuildWebSocketFrame(opcode byte, data []byte, fin bool, mask bool) ([]byte, error) {
	// 检查操作码是否有效
	if !isValidOpCode(opcode) {
		return nil, fmt.Errorf("无效的操作码: %d", opcode)
	}

	// 计算负载长度
	payloadLen := len(data)

	// 计算头部大小
	headerSize := 2 // 基本头部大小
	if payloadLen >= 126 && payloadLen <= 0xFFFF {
		headerSize += 2 // 16位扩展长度
	} else if payloadLen > 0xFFFF {
		headerSize += 8 // 64位扩展长度
	}

	if mask {
		headerSize += 4 // 掩码密钥大小
	}

	// 创建帧缓冲区
	frame := make([]byte, headerSize+payloadLen)

	// 设置第一个字节 (FIN, RSV1, RSV2, RSV3, OpCode)
	if fin {
		frame[0] = 0x80 // 设置FIN位
	}
	frame[0] |= opcode & 0x0F // 设置操作码

	// 设置第二个字节 (Mask, Payload Length)
	if mask {
		frame[1] = 0x80 // 设置掩码位
	}

	// 设置负载长度
	switch {
	case payloadLen < 126:
		frame[1] |= byte(payloadLen)
	case payloadLen <= 0xFFFF:
		frame[1] |= 126 // 16位扩展长度
		binary.BigEndian.PutUint16(frame[2:4], uint16(payloadLen))
	default:
		frame[1] |= 127 // 64位扩展长度
		binary.BigEndian.PutUint64(frame[2:10], uint64(payloadLen))
	}

	// 设置掩码密钥和应用掩码（如果需要）
	if mask {
		maskKey := generateMaskKey()
		maskOffset := 2

		// 如果有扩展长度，调整掩码密钥的偏移量
		if payloadLen >= 126 && payloadLen <= 0xFFFF {
			maskOffset += 2
		} else if payloadLen > 0xFFFF {
			maskOffset += 8
		}

		// 复制掩码密钥到帧
		copy(frame[maskOffset:maskOffset+4], maskKey[:])

		// 应用掩码到负载数据
		framePayloadOffset := maskOffset + 4
		for i := range data {
			frame[framePayloadOffset+i] = data[i] ^ maskKey[i%4]
		}
	} else {
		// 不需要掩码，直接复制数据
		framePayloadOffset := headerSize
		copy(frame[framePayloadOffset:], data)
	}

	return frame, nil
}

var maskRand = rand.Reader

// 生成随机掩码密钥
func generateMaskKey() [4]byte {
	var k [4]byte
	_, _ = io.ReadFull(maskRand, k[:])
	return k
}

// 检查操作码是否有效
func isValidOpCode(opcode byte) bool {
	switch opcode {
	case FrameContinuation, FrameText, FrameBinary,
		FrameClose, FramePing, FramePong:
		return true
	default:
		return false
	}
}

type WSPack struct {
	Header *WebSocketFrameHeader
	Data   []byte
	T      int
}

func (m *WSPack) GetHeader(header interface{}) {
	header = m.Header
}

func (m *WSPack) GetData() []byte {
	return m.Data
}

func (m *WSPack) GetType() int {
	return m.T
}

type WSDataHandle struct {
	IDataHandle
}

func (m *WSDataHandle) Unpacking(ctx context.Context, data []byte) (int, IPack, error) {
	len, header, buf, err := ParseWebSocketFrame(data)
	if err != nil {
		return 0, nil, err
	}
	if buf == nil { //数据不够解析
		return 0, nil, nil
	}

	return len, &WSPack{
		Header: header,
		Data:   buf,
		T:      int(header.OpCode),
	}, nil
}

func (m *WSDataHandle) Packet(ctx context.Context, data []byte) (int, []byte, error) {
	msgType, ok := ctx.Value(WEBSOCKET_CONTEXT_TYPE_KEY).(int)
	if !ok {
		msgType = FrameBinary
	}
	buf, err := BuildWebSocketFrame(byte(msgType), data, true, false) //服务器下发数据不允许使用mask否则大多数浏览器会直接报错
	if err != nil {
		return -1, nil, err
	}
	return len(data), buf, nil
}
