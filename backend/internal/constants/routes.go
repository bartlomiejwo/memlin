package constants

var AuthRoutes = struct {
	GoogleLogin    string
	GoogleCallback string
	RefreshToken   string
}{
	GoogleLogin:    "/api/auth/google",
	GoogleCallback: "/api/auth/callback/google",
}

var WordsRoutes = struct {
	Words string
}{
	Words: "/api/words",
}
