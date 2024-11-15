package persistence_data

import "Eulogist/core/minecraft/protocol/login"

// LoginDataGeneral 是一个通用的结构体，
// 适用于描述 Minecraft 客户端或网易账户的登录数据
type LoginDataGeneral struct {
	IdentityData *login.IdentityData // 网易账户的身份证明
	ClientData   *login.ClientData   // 网易账户的客户端数据
}

// 记录 Minecraft 客户端和网易账户的登录数据
type LoginData struct {
	Client          LoginDataGeneral // 来自 Minecraft 客户端的登录数据
	Server          LoginDataGeneral // 来自网易账户的登录数据
	PlayerUniqueID  int64            // 当前网易账户在当前租赁服所对应的唯一 ID
	PlayerRuntimeID uint64           // 当前网易账户在当前租赁服所对应的运行时 ID
}
