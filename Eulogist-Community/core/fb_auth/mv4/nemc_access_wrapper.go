package fbauth

import (
	fb_client "Eulogist/core/fb_auth/mv4/client"
	"context"
	"encoding/base64"
)

// ...
type AccessWrapper struct {
	ServerCode     string
	ServerPassword string
	Token          string
	Username       string
	Password       string
	Client         *fb_client.Client
}

// ...
func NewAccessWrapper(Client *fb_client.Client, ServerCode, ServerPassword, Token, username, password string) *AccessWrapper {
	return &AccessWrapper{
		Client:         Client,
		ServerCode:     ServerCode,
		ServerPassword: ServerPassword,
		Token:          Token,
		Username:       username,
		Password:       password,
	}
}

// ...
func (aw *AccessWrapper) GetAccess(ctx context.Context, publicKey []byte) (authResponse *fb_client.AuthResponse, err error) {
	pubKeyData := base64.StdEncoding.EncodeToString(publicKey)
	authResponse, err = aw.Client.Auth(ctx, aw.ServerCode, aw.ServerPassword, pubKeyData, aw.Token, aw.Username, aw.Password)
	if err != nil {
		return nil, err
	}
	return
}
