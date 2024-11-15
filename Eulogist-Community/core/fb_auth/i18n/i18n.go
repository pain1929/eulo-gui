package I18n

const (
	Auth_BackendError = iota + 5
	Auth_FailedToRequestEntry
	Auth_HelperNotCreated
	Auth_InvalidFBVersion
	Auth_InvalidHelperUsername
	Auth_InvalidToken
	Auth_InvalidUser
	Auth_ServerNotFound
	Auth_UnauthorizedRentalServerNumber
	Auth_UserCombined
	Auth_FailedToRequestEntry_TryAgain
	Auth_MessageFromAuthServer
)

func T(code uint16) string {
	r, has := I18nDict_zh_CN[code]
	if !has {
		return "???"
	}
	return r
}
