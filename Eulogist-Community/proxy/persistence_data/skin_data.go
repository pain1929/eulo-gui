package persistence_data

import (
	"Eulogist/core/minecraft/protocol"
	"Eulogist/core/tools/skin_process"
)

// 用户的皮肤数据
type SkinData struct {
	NeteaseSkin *skin_process.Skin // 赞颂者处理的皮肤结果
	ServerSkin  *protocol.Skin     // 租赁服返回的最终皮肤信息
}
