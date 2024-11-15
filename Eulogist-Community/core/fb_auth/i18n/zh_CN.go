package I18n

var I18nDict_zh_CN map[uint16]string = map[uint16]string{
	Auth_BackendError:                   "后端错误",
	Auth_FailedToRequestEntry:           "未能请求租赁服入口，请检查租赁服等级设置是否关闭及租赁服密码是否正确。",
	Auth_HelperNotCreated:               "辅助用户尚未创建，请前往用户中心进行创建。",
	Auth_InvalidFBVersion:               "FastBuilder 版本无效，请更新。",
	Auth_InvalidHelperUsername:          "辅助用户的用户名无效，请前往用户中心进行设置。",
	Auth_InvalidToken:                   "无效Token，请重新登录。",
	Auth_InvalidUser:                    "无效用户，请重新登录。",
	Auth_ServerNotFound:                 "租赁服未找到，请检查租赁服是否对所有人开放。",
	Auth_UnauthorizedRentalServerNumber: "对应租赁服号尚未授权，请前往用户中心进行授权。",
	Auth_UserCombined:                   "该用户已经合并到另一个账户中，请使用新账户登录。",
	Auth_FailedToRequestEntry_TryAgain:  "未能请求租赁服入口，请稍后再试。",
	Auth_MessageFromAuthServer:          "来自验证服务器的消息:",
}
