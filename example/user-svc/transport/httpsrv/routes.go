package httpsrv

func routes() []route {
	return []route{
		{
			"SignUp",
			"POST",
			"/signup",
			postSignUpHandler,
		},
		{
			"LogIn",
			"POST",
			"/login",
			postLogInHandler,
		},
		{
			"User",
			"GET",
			"/user",
			getUserHandler,
		},
		{
			"PageUsers",
			"POST",
			"/pageusers",
			postPageUsersHandler,
		},
	}
}
