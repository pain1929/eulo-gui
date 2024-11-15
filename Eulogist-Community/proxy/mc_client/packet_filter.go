package mc_client

import (
	"Eulogist/core/minecraft/protocol/packet"
	raknet_wrapper "Eulogist/core/raknet/wrapper"
	"encoding/json"
	"fmt"
)

/*
数据包过滤器过滤来自 Minecraft 客户端的多个数据包，
然后并将过滤后的多个数据包抄送至网易租赁服。

如果需要，
将根据实际情况由本处的桥接直接发送回应。

writePacketsToServer 指代
用于向客户端抄送数据包的函数。

syncFunc 用于将数据同步到网易租赁服，
它会在 packets 全部被处理完毕后执行，
随后，相应的数据包会被抄送至客户端。

返回的 []error 是一个列表，
分别对应 packets 中每一个数据包的处理成功情况
*/
func (m *MinecraftClient) FiltePacketsAndSendCopy(
	packets []raknet_wrapper.MinecraftPacket,
	writePacketsToServer func(packets []raknet_wrapper.MinecraftPacket),
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
		case *packet.NeteaseJson:
			// 解码 pk.Data 为 JSON 格式
			var jsonMap map[string]any
			err = json.Unmarshal(pk.Data, &jsonMap)
			if err != nil {
				err = fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
				break
			}
			// Login UID 已由赞颂者在先前发送，
			// 所以此处不必重复发送
			if eventName, ok := jsonMap["eventName"].(string); ok {
				if eventName == "LOGIN_UID" {
					shouldSendCopy = false
					break
				}
			}
			// 将 NetEase UID 修正为真实值
			if _, ok := jsonMap["uid"]; ok {
				jsonMap["uid"] = fmt.Sprintf("%d", m.PersistenceData.LoginData.Server.IdentityData.Uid)
			}
			// 将 JSON 重新编码到 pk.Data
			pk.Data, err = json.Marshal(jsonMap)
			if err != nil {
				err = fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
				break
			}
			// 要求该数据包需要经编码后发送
			minecraftPacket.Bytes = nil
		default:
			// 默认情况下，我们需要将
			// 数据包同步到网易租赁服
		}
		// 提交子结果
		errResults = append(errResults, err)
		if shouldSendCopy {
			sendCopy = append(sendCopy, minecraftPacket)
		}
	}
	// 同步数据并抄送数据包
	err := syncFunc()
	writePacketsToServer(sendCopy)
	// 返回值
	if err != nil {
		return errResults, fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
	} else {
		return errResults, nil
	}
}
