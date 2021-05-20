package httpsrv

import (
	"net/http"
)

type UserServiceHandler interface {
	PostSignUp(w http.ResponseWriter, r *http.Request)
	PostLogIn(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	PostPageUsers(w http.ResponseWriter, r *http.Request)
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func routes(handler UserServiceHandler) []route {
	return []route{
		{
			"SignUp",
			"POST",
			"/userservice/signup",
			handler.PostSignUp,
		},
		{
			"LogIn",
			"POST",
			"/userservice/login",
			handler.PostLogIn,
		},
		{
			"User",
			"GET",
			"/userservice/user",
			handler.GetUser,
		},
		{
			"PageUsers",
			"POST",
			"/userservice/pageusers",
			handler.PostPageUsers,
		},
	}
}
