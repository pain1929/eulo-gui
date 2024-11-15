package mc_server

import (
	fbauth "Eulogist/core/fb_auth/mv4"
	"Eulogist/core/minecraft/protocol/packet"
	raknet_wrapper "Eulogist/core/raknet/wrapper"
	"Eulogist/core/tools/py_rpc"
	"encoding/json"
	"fmt"
)

// OnPyRpc 处理数据包 PyRpc
func (m *MinecraftServer) OnPyRpc(p *packet.PyRpc) (shouldSendCopy bool, err error) {
	// 解码 PyRpc
	if p.Value == nil {
		return true, nil
	}
	content, err := py_rpc.Unmarshal(p.Value)
	if err != nil {
		return true, fmt.Errorf("OnPyRpc: %v", err)
	}
	// 根据内容类型处理不同的 PyRpc
	switch c := content.(type) {
	case *py_rpc.StartType:
		c.Content = fbauth.TransferData(m.fbClient, c.Content)
		c.Type = py_rpc.StartTypeResponse
		m.Conn.WriteSinglePacket(
			raknet_wrapper.MinecraftPacket{
				Packet: &packet.PyRpc{
					Value:         py_rpc.Marshal(c),
					OperationType: packet.PyRpcOperationTypeSend,
				},
			},
		)
	case *py_rpc.GetMCPCheckNum:
		// 如果已完成零知识证明(挑战)，
		// 则不做任何操作
		if m.getCheckNumEverPassed {
			break
		}
		// 创建请求并发送到认证服务器并获取响应
		arg, _ := json.Marshal([]any{
			c.FirstArg,
			c.SecondArg.Arg,
			m.PersistenceData.LoginData.PlayerUniqueID,
		})
		ret := fbauth.TransferCheckNum(m.fbClient, string(arg))
		// 解码响应并调整数据
		ret_p := []any{}
		json.Unmarshal([]byte(ret), &ret_p)
		if len(ret_p) > 7 {
			ret6, ok := ret_p[6].(float64)
			if ok {
				ret_p[6] = int64(ret6)
			}
		}
		// 完成零知识证明(挑战)
		m.Conn.WriteSinglePacket(
			raknet_wrapper.MinecraftPacket{
				Packet: &packet.PyRpc{
					Value:         py_rpc.Marshal(&py_rpc.SetMCPCheckNum{ret_p}),
					OperationType: packet.PyRpcOperationTypeSend,
				},
			},
		)
		// 标记零知识证明(挑战)已在当前会话下永久完成
		m.getCheckNumEverPassed = true
	default:
		// 对于其他种类的 PyRpc 数据包，
		// 返回 true 表示需要将数据包抄送至
		// Minecraft 客户端
		return true, nil
	}
	// 返回值
	return false, nil
}
