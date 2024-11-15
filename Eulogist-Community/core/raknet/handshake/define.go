package handshake

// saltClaims 是保存服务器在
// ServerToClientHandshake
// 数据包中发送的 salt 的声明
type saltClaims struct {
	Salt string `json:"salt"`
}
