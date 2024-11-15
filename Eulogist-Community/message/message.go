package message

import (
	"encoding/binary"
	"fmt"
	"net"
)

// MessageHdr 头部结构体
type MessageHdr struct {
	Zero   uint16 // 头部，固定为 0xAB
	MsgLen uint16 // 消息长度 不包含头部
}

// NormalMsg 结构体
type NormalMsg struct {
	Started     bool    // 端口启动成功
	ErrorMsgLen uint16  // 错误消息长度
	ErrorMsg    []uint8 // 错误消息
}

// setMsg 方法：根据参数设置 NormalMsg 内容
func (msg *NormalMsg) setMsg(started bool, errorMsg string) {
	// 设置 Started 和 ErrorMsgLen
	msg.Started = started
	msg.ErrorMsg = []uint8(errorMsg)

	// 设置 ErrorMsgLen
	msg.ErrorMsgLen = uint16(len(msg.ErrorMsg))
}

// 将 NormalMsg 转换为字节切片
func (msg *NormalMsg) toBytes() []uint8 {
	// 计算消息总长度：Started (1 字节) + ErrorMsgLen (2 字节) + 错误消息 (ErrorMsgLen 字节)
	totalLen := 1 + 2 + len(msg.ErrorMsg)

	// 创建字节切片
	buf := make([]uint8, totalLen)

	// 将 Started 转换为 uint8 (1 字节)
	if msg.Started {
		buf[0] = 1
	} else {
		buf[0] = 0
	}

	// 将 ErrorMsgLen 转换为 uint16 (2 字节)
	binary.LittleEndian.PutUint16(buf[1:3], msg.ErrorMsgLen)

	// 将 ErrorMsg 复制到字节切片
	copy(buf[3:], msg.ErrorMsg)

	return buf
}

// sendMsg 方法：根据参数构建 MessageHdr 和 NormalMsg 并通过 TCP 连接发送
func SendMsg(started bool, errorMsg string, conn net.Conn) error {
	// 创建 NormalMsg 并设置内容
	var msg NormalMsg
	msg.setMsg(started, errorMsg)

	// 获取 NormalMsg 的字节表示
	msgBytes := msg.toBytes()

	// 设置消息头的内容
	var header MessageHdr
	header.Zero = 0xAB
	header.MsgLen = uint16(len(msgBytes)) // 仅包含 NormalMsg 的长度，不包括头部长度

	// 将头部转换为字节切片
	headerBytes := make([]uint8, 4)
	binary.LittleEndian.PutUint16(headerBytes[0:2], header.Zero)
	binary.LittleEndian.PutUint16(headerBytes[2:4], header.MsgLen)

	// 发送头部
	_, err := conn.Write(headerBytes)
	if err != nil {
		return fmt.Errorf("failed to send header: %v", err)
	}

	// 发送消息内容
	_, err = conn.Write(msgBytes)
	if err != nil {
		return fmt.Errorf("failed to send message body: %v", err)
	}

	return nil
}
