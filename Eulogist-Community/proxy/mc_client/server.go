package mc_client

import (
	raknet_wrapper "Eulogist/core/raknet/wrapper"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"

	"Eulogist/core/minecraft/raknet"
)

// CreateListener 在 127.0.0.1 上以 Raknet 协议侦听 Minecraft 客户端的连接，
// 这意味着您成功创建了一个 Minecraft 数据包代理服务器。
// 稍后，您可以通过 m.GetServerAddress 来取得服务器地址
func (m *MinecraftClient) CreateListener() error {
	// 创建一个 Raknet 监听器
	listener, err := raknet.Listen("127.0.0.1:1929")
	if err != nil {
		return fmt.Errorf("CreateListener: %v", err)
	}
	// 获取监听器的地址
	address, ok := listener.Addr().(*net.UDPAddr)
	if !ok {
		return fmt.Errorf("CreateListener: Failed to get address for listener")
	}
	// 初始化变量
	m.listener = listener
	m.connected = make(chan struct{}, 1)
	m.Address = address
	m.Conn = raknet_wrapper.NewRaknet()
	// 返回成功
	return nil
}

// WaitConnect 等待 Minecraft 客户端连接到服务器
func (m *MinecraftClient) WaitConnect() error {
	// 接受客户端连接
	conn, err := m.listener.Accept()
	if err != nil {
		return fmt.Errorf("WaitConnect: %v", err)
	}
	// 丢弃其他连接
	go func() {
		for {
			conn, err := m.listener.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	// 初始化变量
	serverKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	m.Conn.SetConnection(conn, serverKey)
	m.connected <- struct{}{}
	// 返回成功
	return nil
}
