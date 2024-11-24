package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// MsgType 消息类型
type MsgType uint64

const (
	EuloMsg    MsgType = 0 // eulogist 进程返回的启动信息
	SendCmdMsg MsgType = 1 // GUI 进程发送的命令信息
)

// MessageHdr 头部结构体
type MessageHdr struct {
	Zero   uint16 // 头部，固定为 0xAB
	Pad    uint16 // 填充
	MsgLen uint32 // 消息长度 (不包含头部，但包含 MsgType 和消息体长度)
}

// 所有的tcp通讯必须使用以下结构
// MessageHdr（8字节） + MsgType（8字节） + 消息结构 例如：EuloMsgType

// EuloMsgType 通知 gui进程启动成功或者失败
type EuloMsgType struct {
	Started     bool    // 端口启动成功
	Pad         int8    // 填充
	ErrorMsgLen uint16  // 错误消息长度
	ErrorMsg    []uint8 // 错误消息
}

// SetMsg 方法：设置是否启动成功
func (msg *EuloMsgType) SetMsg(started bool, errorMsg string) {
	// 设置 Started 和 ErrorMsgLen
	msg.Started = started
	msg.ErrorMsg = []uint8(errorMsg)

	// 设置 ErrorMsgLen
	msg.ErrorMsgLen = uint16(len(msg.ErrorMsg))
}

// toBytes 序列化 MessageHdr 为字节切片，返回切片及其长度
func (hdr *MessageHdr) toBytes() ([]byte, int, error) {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, hdr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialize MessageHdr: %v", err)
	}
	data := buffer.Bytes()
	return data, len(data), nil
}

// toBytes 序列化 EuloMsgType 为字节切片，返回切片及其长度
func (msg *EuloMsgType) toBytes() ([]byte, int, error) {
	var buffer bytes.Buffer

	// 写入 Started
	err := binary.Write(&buffer, binary.LittleEndian, msg.Started)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialize Started: %v", err)
	}

	// 写入 Pad
	err = binary.Write(&buffer, binary.LittleEndian, msg.Pad)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialize Pad: %v", err)
	}

	// 写入 ErrorMsgLen
	err = binary.Write(&buffer, binary.LittleEndian, msg.ErrorMsgLen)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialize ErrorMsgLen: %v", err)
	}

	// 写入 ErrorMsg
	err = binary.Write(&buffer, binary.LittleEndian, msg.ErrorMsg)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to serialize ErrorMsg: %v", err)
	}

	data := buffer.Bytes()
	return data, len(data), nil
}

// SendPacket 封装 EuloMsgType 的发送逻辑
func (msg *EuloMsgType) SendPacket(conn net.Conn) error {
	// 获取消息内容
	messageBody, bodyLen, err := msg.toBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize EuloMsgType: %v", err)
	}

	// 构造消息头
	hdr := &MessageHdr{
		Zero:   0xAB,                                              // 固定值
		Pad:    0x00,                                              // 填充
		MsgLen: uint32(bodyLen) + uint32(binary.Size(MsgType(0))), // 消息类型 + 消息体长度
	}

	// 获取消息头字节
	hdrBytes, _, err := hdr.toBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize MessageHdr: %v", err)
	}

	// 将消息头 + 消息类型 + 消息体拼接到一个缓冲区
	var buffer bytes.Buffer

	// 写入消息头
	buffer.Write(hdrBytes)

	// 写入消息类型
	err = binary.Write(&buffer, binary.LittleEndian, MsgType(EuloMsg))
	if err != nil {
		return fmt.Errorf("failed to write MsgType: %v", err)
	}

	// 写入消息体
	buffer.Write(messageBody)

	// 将完整消息发送到 TCP 连接
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
