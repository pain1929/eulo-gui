package mc_server

import (
	"Eulogist/core/minecraft/protocol/packet"
	"Eulogist/core/raknet/handshake"
	raknet_wrapper "Eulogist/core/raknet/wrapper"
	"Eulogist/core/tools/py_rpc"
	"fmt"
)

/*
数据包过滤器过滤来自租赁服的多个数据包，
然后并将过滤后的多个数据包抄送至客户端。

如果需要，
将根据实际情况由本处的桥接直接发送回应。

writeSinglePacketToClient 指代
用于向客户端抄送数据包的函数。

syncFunc 用于将数据同步到 Minecraft，
它会在 packets 全部被处理完毕后执行，
随后，相应的数据包会被抄送至网易租赁服。

返回的 []error 是一个列表，
分别对应 packets 中每一个数据包的处理成功情况
*/
func (m *MinecraftServer) FiltePacketsAndSendCopy(
	packets []raknet_wrapper.MinecraftPacket,
	writePacketsToClient func(packets []raknet_wrapper.MinecraftPacket),
	syncFunc func() error,
) (errResults []error, syncError error) {
	// 初始化
	errResults = make([]error, 0)
	sendCopy := make([]raknet_wrapper.MinecraftPacket, 0)
	// 处理每个数据包
	for _, minecraftPacket := range packets {
		// 初始化
		var shouldSendCopy bool = true
		var err error
		// 根据数据包的类型进行不同的处理
		switch pk := minecraftPacket.Packet.(type) {
		case *packet.PyRpc:
			shouldSendCopy, err = m.OnPyRpc(pk)
			if err != nil {
				err = fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
			}
		case *packet.StartGame:
			// 预处理
			m.PersistenceData.LoginData.PlayerUniqueID, m.PersistenceData.LoginData.PlayerRuntimeID = handshake.HandleStartGame(m.Conn, pk)
			playerSkin := m.PersistenceData.SkinData.NeteaseSkin
			// 发送简要身份证明
			m.Conn.WriteSinglePacket(raknet_wrapper.MinecraftPacket{
				Packet: &packet.NeteaseJson{
					Data: []byte(
						fmt.Sprintf(
							`{"eventName":"LOGIN_UID","resid":"","uid":"%d"}`,
							m.PersistenceData.LoginData.Server.IdentityData.Uid,
						),
					),
				},
			})
			// 其他组件处理
			if playerSkin == nil {
				m.Conn.WriteSinglePacket(raknet_wrapper.MinecraftPacket{
					Packet: &packet.PyRpc{
						Value:         py_rpc.Marshal(&py_rpc.SyncUsingMod{}),
						OperationType: packet.PyRpcOperationTypeSend,
					},
				})
			} else {
				// 初始化
				modUUIDs := make([]any, 0)
				outfitInfo := make(map[string]int64, 0)
				// 设置数据
				for modUUID, outfitType := range m.PersistenceData.BotComponent {
					modUUIDs = append(modUUIDs, modUUID)
					if outfitType != nil {
						outfitInfo[modUUID] = int64(*outfitType)
					}
				}
				// 组件处理
				m.Conn.WriteSinglePacket(raknet_wrapper.MinecraftPacket{
					Packet: &packet.PyRpc{
						Value: py_rpc.Marshal(&py_rpc.SyncUsingMod{
							modUUIDs,
							playerSkin.SkinUUID,
							playerSkin.SkinItemID,
							true,
							outfitInfo,
						}),
						OperationType: packet.PyRpcOperationTypeSend,
					},
				})
			}
		case *packet.UpdatePlayerGameType:
			if pk.PlayerUniqueID == m.PersistenceData.LoginData.PlayerUniqueID {
				// 如果玩家的唯一 ID 与数据包中记录的值匹配，
				// 则向客户端发送 SetPlayerGameType 数据包，
				// 并放弃当前数据包的发送，
				// 以确保 Minecraft 客户端可以正常同步游戏模式更改。
				// 否则，按原样抄送当前数据包
				sendCopy = append(sendCopy, raknet_wrapper.MinecraftPacket{
					Packet: &packet.SetPlayerGameType{GameType: pk.GameType},
				})
				shouldSendCopy = false
			}
		default:
			// 默认情况下，
			// 我们需要将数据包同步到客户端
		}
		// 提交子结果
		errResults = append(errResults, err)
		if shouldSendCopy {
			sendCopy = append(sendCopy, minecraftPacket)
		}
	}
	// 同步数据并抄送数据包
	err := syncFunc()
	writePacketsToClient(sendCopy)
	// 返回值
	if err != nil {
		return errResults, fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
	} else {
		return errResults, nil
	}
}
