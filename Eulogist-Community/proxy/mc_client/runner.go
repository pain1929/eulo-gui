package mc_client

import (
	"Eulogist/core/minecraft/protocol/packet"
	"Eulogist/core/raknet/handshake"
	"Eulogist/proxy/persistence_data"
	"fmt"
)

/*
RunServer 在本机上运行一个代理服务器以等待 Minecraft 客户端连接。
persistenceData 用于设置持久化数据，
它应该与 MinecraftServer 使用同一个。

服务器开放的端口将被自动设置，
您可以使用 client.Address.Port 来取得开放的端口。

当 Minecraft 连接时，管道 connected 将收到数据
*/
func RunServer(persistenceData *persistence_data.PersistenceData) (
	client *MinecraftClient,
	connected chan struct{},
	err error,
) {
	// prepare
	client = &MinecraftClient{PersistenceData: persistenceData}
	// start listening
	err = client.CreateListener()
	if err != nil {
		return nil, nil, fmt.Errorf("RunServer: %v", err)
	}
	// wait Minecraft Client to connect
	go func() {
		err = client.WaitConnect()
		if err != nil {
			panic(fmt.Sprintf("RunServer: %v", err))
		}
		client.Conn.ProcessIncomingPackets()
	}()
	// return
	return client, client.connected, nil
}

// WaitClientHandshakeDown 等待 Minecraft
// 完成与赞颂者的基本数据包交换。
// 此函数应当只被调用一次
func (m *MinecraftClient) WaitClientHandshakeDown() error {
	// 准备
	var err error
	// 处理来自 Minecraft
	// 客户端的登录相关的数据包
	for {
		for _, pk := range m.Conn.ReadPackets() {
			// 处理初始连接数据包
			switch p := pk.Packet.(type) {
			case *packet.RequestNetworkSettings:
				err = handshake.HandleRequestNetworkSettings(m.Conn, p)
				if err != nil {
					panic(fmt.Sprintf("WaitClientHandshakeDown: %v", err))
				}
			case *packet.Login:
				m.PersistenceData.LoginData.Client.IdentityData, m.PersistenceData.LoginData.Client.ClientData, err = handshake.HandleLogin(m.Conn, p)
				if err != nil {
					panic(fmt.Sprintf("WaitClientHandshakeDown: %v", err))
				}
			case *packet.ClientToServerHandshake:
				// 连接已完成初始化，
				// 于是我们返回值
				return nil
			}
		}
		// 检查连接状态
		select {
		case <-m.Conn.GetContext().Done():
			return fmt.Errorf("WaitClientHandshakeDown: Minecraft closed its connection to eulogist")
		default:
		}
	}
}
