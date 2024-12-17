package Eulogist

// NetEaseConfig 结构体定义了一个可以允许
// Minecraft 客户端自动进入特定服务器地址的配置文件字段
type NetEaseConfig struct {
	RoomInfo   RoomInfo   `json:"room_info"`   // 房间信息
	PlayerInfo PlayerInfo `json:"player_info"` // 玩家信息
	SkinInfo   SkinInfo   `json:"skin_info"`   // 皮肤信息
	Misc       Misc       `json:"misc"`        // 杂项信息
}

// 房间信息
type RoomInfo struct {
	IP   string `json:"ip"`   // 房间的IP地址
	Port int    `json:"port"` // 房间的端口号
}

// 玩家信息
type PlayerInfo struct {
	UserID   int    `json:"user_id"`   // 玩家的用户 ID
	UserName string `json:"user_name"` // 玩家的用户名
	Urs      string `json:"urs"`       // 玩家的 URS
}

// 皮肤信息
type SkinInfo struct {
	// 皮肤的路径。
	// 对于普通皮肤，这指向一个 PNG 文件，
	// 对于高级皮肤(如 4D 皮肤)，
	// 这指向一个 ZIP 压缩包
	SkinPath string `json:"skin"`
	// 描述皮肤的手臂是否纤细
	Slim bool `json:"slim"`
}

// 杂项信息
type Misc struct {
	// 多人游戏类型
	MultiplayerGameType int `json:"multiplayer_game_type"`
}

// DefaultNetEaseConfig 函数返回一个默认的 NetEaseConfig 实例
func DefaultNetEaseConfig() NetEaseConfig {
	return NetEaseConfig{
		RoomInfo: RoomInfo{IP: "127.0.0.1", Port: 19132},
		PlayerInfo: PlayerInfo{
			UserID:   -1,
			UserName: "可可爱潇净",
			Urs:      "***",
		},
		Misc: Misc{MultiplayerGameType: 100},
	}
}
