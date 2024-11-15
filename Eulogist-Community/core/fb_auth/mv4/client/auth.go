package fb_client

import (
	I18n "Eulogist/core/fb_auth/i18n"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// 客户端向验证服务器发送的请求体，
// 用于获得 FBToken，
// 或使客户端登录到网易租赁服。
// AuthResponse 是该请求体对应的响应体
type AuthRequest struct {
	/*
		此字段非空时，则下方 UserName 和 Password 为空，
		否则反之。

		当 FBToken 或 UserName、Password 二者中任意一个
		填写的值正确时，用户将登录到用户中心，然后进入租赁服
	*/
	FBToken string `json:"login_token,omitempty"`

	UserName string `json:"username,omitempty"` // 用户在用户中心的用户名
	Password string `json:"password,omitempty"` // 用户在用户中心的密码

	ServerCode     string `json:"server_code"`     // 要进入的租赁服的 服务器号
	ServerPassword string `json:"server_passcode"` // 该租赁服的 密码

	ClientPublicKey string `json:"client_public_key"` // ...
}

// 验证服务器对 AuthRequest 的响应体
type AuthResponse struct {
	/*
		描述请求的成功状态。

		如果成功，则其余的所有非可选字段都将有值，
		这也包括 Message 本身。

		如果失败，除了本字段和 Message 以外，
		其余所有字段都为默认的零值，
		同时 Message 会展示对应的失败原因
	*/
	SuccessStates bool   `json:"success"`
	ServerMessage string `json:"server_msg,omitempty"` // 来自验证服务器的消息
	Message

	BotLevel     int             `json:"growth_level"`          // 机器人的等级
	BotSkin      SkinInfo        `json:"skin_info,omitempty"`   // 机器人的皮肤信息
	BotComponent map[string]*int `json:"outfit_info,omitempty"` // 机器人当前已加载的组件及其附加值

	FBToken    string `json:"token"`      // ...
	MasterName string `json:"respond_to"` // 机器人主人的游戏名称

	RentalServerIP string `json:"ip_address"` // 欲登录的租赁服的 IP 地址
	ChainInfo      string `json:"chainInfo"`  // 欲登录的租赁服的链请求
}

// 描述 AuthResponse 所附带的额外信息
type Message struct {
	/*
		若 AuthRequest 成功，
		则对于原生的 FastBuilder 的验证服务器(mv4)，
		此字段为 "正常返回"；
		否则，对于 咕咕酱及其开发团队 的验证服务器，
		此字段为 "well down"。

		当 AuthRequest 失败时，
		若此字段非空，则它将阐明对应的失败原因，
		否则，由下方的 Translation 揭示具体的原因
	*/
	Information string `json:"message,omitempty"`
	// 表示错误码，且可以与 i18n 中所记的映射对应。
	// 如果不存在，则其默认值为 0，
	// 如果未使用，则其默认值为 -1
	Translation int `json:"translation,omitempty"`
}

// 描述 AuthResponse 中附带的 机器人 的皮肤信息
type SkinInfo struct {
	ItemID          string `json:"entity_id"` // 皮肤的资源 ID
	SkinDownloadURL string `json:"res_url"`   // 皮肤的下载地址 [需要验证]
	SkinIsSlim      bool   `json:"is_slim"`   // 皮肤的手臂是否纤细
}

// Ret: chain, ip, token, error
func (client *Client) Auth(
	ctx context.Context,
	serverCode string, serverPassword string,
	key string,
	fbtoken string, username string, password string,
) (*AuthResponse, error) {
	var authResponse AuthResponse
	// prepare
	request := AuthRequest{
		FBToken:         fbtoken,
		UserName:        username,
		Password:        password,
		ServerCode:      serverCode,
		ServerPassword:  serverPassword,
		ClientPublicKey: key,
	}
	authRequest, _ := json.Marshal(request)
	// pack request and marshal to binary
	httpResponse, err := client.HttpClient.Post(
		fmt.Sprintf("%s/api/phoenix/login", client.AuthServer),
		"application/json",
		bytes.NewBuffer(authRequest),
	)
	if err != nil {
		panic(fmt.Sprintf("Auth: %v", err))
	}
	authResponse = AssertAndParse[AuthResponse](httpResponse)
	// get response
	if !authResponse.SuccessStates {
		failedReason := authResponse.Message
		err := failedReason.Information
		if t := failedReason.Translation; t != -1 && t != 0 {
			err = I18n.T(uint16(t))
		}
		return nil, fmt.Errorf("%s", err)
	}
	// if reuqest failed
	client.ClientInfo.GrowthLevel = authResponse.BotLevel
	// set value
	return &authResponse, nil
	// return
}
