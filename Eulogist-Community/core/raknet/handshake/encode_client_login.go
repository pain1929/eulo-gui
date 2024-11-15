package handshake

import (
	fb_client "Eulogist/core/fb_auth/mv4/client"
	"Eulogist/core/minecraft/protocol"
	"Eulogist/core/minecraft/protocol/login"
	"Eulogist/core/tools/skin_process"
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// EncodeLogin 编码登录请求。
// 它使用提供的身份验证响应、
// 客户端密钥和皮肤信息生成登录请求数据包
func EncodeLogin(
	authResponse *fb_client.AuthResponse,
	clientKey *ecdsa.PrivateKey,
	skin *skin_process.Skin,
) (
	request []byte,
	identityData *login.IdentityData, clientData *login.ClientData,
	err error,
) {
	identity := login.IdentityData{}
	client := login.ClientData{}

	// 设置默认的身份和客户端数据
	defaultIdentityData(&identity)
	defaultClientData(&client, authResponse, skin)

	// 我们以 Android 设备登录，这将在 JWT 链中的 titleId 字段中显示。
	// 这些字段无法被编辑，而我们也仅仅是强制以 Android 数据进行登录
	setAndroidData(&client)

	// 编码登录请求
	request = login.Encode(authResponse.ChainInfo, client, clientKey)
	// 解析身份数据以确保其有效
	identity, client, _, err = login.Parse(request)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("EncodeLogin: WARNING: Identity data parsing error: %v", err)
	}

	return request, &identity, &client, nil
}

// defaultIdentityData 编辑传入的 IdentityData，
// 为所有未更改的字段设置默认值
func defaultIdentityData(data *login.IdentityData) {
	if data.Identity == "" {
		data.Identity = uuid.New().String()
	}
	if data.DisplayName == "" {
		data.DisplayName = "Steve"
	}
}

// PhoenixBuilder specific changes.
// Author: Happy2018new
//
// defaultClientData 编辑传入的 ClientData，
// 为所有未更改的字段设置默认值
func defaultClientData(
	d *login.ClientData,
	authResponse *fb_client.AuthResponse,
	skin *skin_process.Skin,
) {
	rand.Seed(time.Now().Unix())

	d.ServerAddress = authResponse.RentalServerIP
	d.DeviceOS = protocol.DeviceAndroid
	d.GameVersion = protocol.CurrentVersion
	d.GrowthLevel = authResponse.BotLevel // Netease
	d.ClientRandomID = rand.Int63()
	d.DeviceID = uuid.New().String()
	d.LanguageCode = "zh_CN" // Netease
	d.AnimatedImageData = make([]login.SkinAnimation, 0)
	d.PersonaPieces = make([]login.PersonaPiece, 0)
	d.PieceTintColours = make([]login.PersonaPieceTintColour, 0)
	d.SelfSignedID = uuid.New().String()

	if skin != nil {
		if len(skin.SkinUUID) == 0 {
			skin.SkinUUID = uuid.NewString()
		}
		d.SkinID = skin.SkinUUID
		d.SkinItemID = skin.SkinItemID
		d.SkinData = base64.StdEncoding.EncodeToString(skin.SkinPixels)
		d.SkinImageHeight, d.SkinImageWidth = skin.SkinHight, skin.SkinWidth
		d.SkinGeometry = base64.StdEncoding.EncodeToString(skin.SkinGeometry)
		d.SkinGeometryVersion = base64.StdEncoding.EncodeToString([]byte("0.0.0"))
		d.SkinResourcePatch = base64.StdEncoding.EncodeToString(skin.SkinResourcePatch)
		d.PremiumSkin = true
	} else {
		d.SkinID = uuid.New().String()
		d.SkinData = base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0, 0, 0, 255}, 32*64))
		d.SkinGeometry = base64.StdEncoding.EncodeToString(skin_process.DefaultSkinGeometry)
		d.SkinResourcePatch = base64.StdEncoding.EncodeToString(skin_process.DefaultWideSkinResourcePatch)
		d.SkinImageHeight = 32
		d.SkinImageWidth = 64
	}
}

// setAndroidData 确保传入的 login.ClientData
// 匹配您在 Android 设备上看到的设置
func setAndroidData(data *login.ClientData) {
	data.DeviceOS = protocol.DeviceAndroid
	data.GameVersion = protocol.CurrentVersion
}
