package Eulogist

// 启动类型常量
const (
	// 默认启动方式，
	// 此方式下将会自动启动 Minecraft 客户端，
	// 并自动使该客户端连接到赞颂者
	LaunchTypeNormal int = iota
	// 旧启动方法，
	// 用户需自行使用 Minecraft 客户端连接到赞颂者
	LaunchTypeNoOperation
)

// 验证服务器地址
const (
	AuthServerAddressFastBuilder = "https://user.fastbuilder.pro"
	AuthServerAddressLiliya233   = "https://liliya233.uk"
)

// EulogistConfig 结构体定义了 Eulogist 的配置信息
type EulogistConfig struct {
	LaunchType int    `json:"launch_type"`       // 启动类型
	NEMCPath   string `json:"nemc_program_path"` // Minecraft 客户端的程序路径

	// 自定义皮肤路径。
	// 目的仅在于在 Minecraft 客户端处使用并显示该皮肤，
	// 但此皮肤数据实际上并没有被同步到租赁服，
	// 因此实际的皮肤表现，请以实际情况为准
	SkinPath string `json:"skin_path"`

	RentalServerCode     string `json:"rental_server_code"`     // 网易租赁服编号
	RentalServerPassword string `json:"rental_server_password"` // 该租赁服对应的密码

	FBToken string `json:"fb_token"`
}

// DefaultEulogistConfig 返回一个默认的 EulogistConfig 实例
func DefaultEulogistConfig() EulogistConfig {
	return EulogistConfig{
		LaunchType: LaunchTypeNormal,
		NEMCPath:   `Minecraft.Windows.exe`,
	}
}

// LookUpAuthServerAddress 根据令牌查找认证服务器地址
func LookUpAuthServerAddress(token string) string {
	if len(token) < 3 {
		return AuthServerAddressFastBuilder
	}

	switch token[:3] {
	case "w9/":
		return AuthServerAddressFastBuilder
	case "y8/":
		return AuthServerAddressLiliya233
	default:
		return AuthServerAddressFastBuilder
	}
}
