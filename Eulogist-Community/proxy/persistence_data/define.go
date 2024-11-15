package persistence_data

// 当前用户的持久化数据
type PersistenceData struct {
	LoginData    LoginData       // Minecraft 客户端和网易账户的登录数据
	SkinData     SkinData        // 用户的皮肤信息
	BotComponent map[string]*int // 用户当前已加载的网易组件(如法阵)及其附加值
}
