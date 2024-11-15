package mc_server

import (
	fbauth "Eulogist/core/fb_auth/mv4"
	fb_client "Eulogist/core/fb_auth/mv4/client"
	"Eulogist/core/minecraft/protocol"
	"Eulogist/core/minecraft/protocol/packet"
	"Eulogist/core/raknet/handshake"
	raknet_wrapper "Eulogist/core/raknet/wrapper"
	"Eulogist/core/tools/skin_process"
	"Eulogist/proxy/persistence_data"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"time"

	"Eulogist/core/minecraft/raknet"
)

// ConnectToServer 用于连接到 basicConfig 所指代的租赁服。
// persistenceData 用于设置持久化数据，
// 它应该与 MinecraftClient 使用同一个
func ConnectToServer(
	basicConfig BasicConfig,
	persistenceData *persistence_data.PersistenceData,
) (*MinecraftServer, error) {
	// 准备
	mcServer := MinecraftServer{
		fbClient:        fb_client.CreateClient(basicConfig.AuthServer),
		PersistenceData: persistenceData,
	}
	// 初始化
	authenticator := fbauth.NewAccessWrapper(
		mcServer.fbClient,
		basicConfig.ServerCode,
		basicConfig.ServerPassword,
		basicConfig.Token,
		"", "",
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	// 向验证服务器请求信息
	clientKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	armoured_key, _ := x509.MarshalPKIXPublicKey(&clientKey.PublicKey)
	authResponse, err := authenticator.GetAccess(ctx, armoured_key)
	if err != nil {
		return nil, fmt.Errorf("ConnectToServer: %v", err)
	}
	// 设置皮肤信息
	if len(authResponse.BotSkin.SkinDownloadURL) > 0 {
		mcServer.PersistenceData.SkinData.NeteaseSkin = new(skin_process.Skin)
		err = skin_process.GetSkinFromAuthResponse(authResponse, mcServer.PersistenceData.SkinData.NeteaseSkin)
		if err != nil {
			return nil, fmt.Errorf("ConnectToServer: %v", err)
		}
	}
	// 连接到服务器
	connection, err := raknet.DialContext(ctx, authResponse.RentalServerIP)
	if err != nil {
		return nil, fmt.Errorf("ConnectToServer: %v", err)
	}
	// 同步数据
	mcServer.PersistenceData.BotComponent = authResponse.BotComponent
	mcServer.authResponse = authResponse
	mcServer.Conn = raknet_wrapper.NewRaknet()
	// 设置底层连接并启动数据包解析
	mcServer.Conn.SetConnection(connection, clientKey)
	go mcServer.Conn.ProcessIncomingPackets()
	// 返回值
	return &mcServer, nil
}

/*
FinishHandshake 用于赞颂者完成
与网易租赁服的基本数据包交换。

在与网易租赁服建立 Raknet 连接后，
由赞颂者发送第一个数据包，
用于向服务器请求网络信息设置。

随后，得到来自网易服务器的回应，
并由赞颂者完成基础登录序列，
然后，最终完成与网易租赁服的握手。

此函数应当只被调用一次
*/
func (m *MinecraftServer) FinishHandshake() error {
	// 准备
	var err error
	// 向网易租赁服请求网络设置，
	// 这是赞颂者登录到网易租赁服的第一个数据包
	m.Conn.WriteSinglePacket(
		raknet_wrapper.MinecraftPacket{
			Packet: &packet.RequestNetworkSettings{ClientProtocol: protocol.CurrentProtocol},
		},
	)
	// 处理来自 bot 端的登录相关数据包
	for {
		for _, pk := range m.Conn.ReadPackets() {
			// 处理初始连接数据包
			switch p := pk.Packet.(type) {
			case *packet.NetworkSettings:
				m.PersistenceData.LoginData.Server.IdentityData, m.PersistenceData.LoginData.Server.ClientData, err = handshake.HandleNetworkSettings(
					m.Conn, p, m.authResponse, m.PersistenceData.SkinData.NeteaseSkin,
				)
				if err != nil {
					return fmt.Errorf("FinishHandshake: %v", err)
				}
			case *packet.ServerToClientHandshake:
				err = handshake.HandleServerToClientHandshake(m.Conn, p)
				if err != nil {
					return fmt.Errorf("FinishHandshake: %v", err)
				}
				// 连接已完成初始化，
				// 于是我们返回值
				return nil
			}
		}
		// 检查连接状态
		select {
		case <-m.Conn.GetContext().Done():
			return fmt.Errorf("FinishHandshake: NetEase Minecraft Rental Server closed their connection to eulogist")
		default:
		}
	}
}
