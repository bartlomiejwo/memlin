package localization

var AuthKeys = struct {
	MissingRedirectURI     string
	MissingAuthCode        string
	AuthFail               string
	UserInfoRetrieveFail   string
	UserInfoParseFail      string
	UserCreateUpdateFail   string
	JWTCreateFail          string
	RefreshTokenCreateFail string
}{
	MissingRedirectURI:     "missing_redirect_uri",
	MissingAuthCode:        "missing_auth_code",
	AuthFail:               "auth_fail",
	UserInfoRetrieveFail:   "user_info_retrieve_fail",
	UserInfoParseFail:      "user_info_parse_fail",
	UserCreateUpdateFail:   "user_create_update_fail",
	JWTCreateFail:          "jwt_create_fail",
	RefreshTokenCreateFail: "refresh_token_create_fail",
}
