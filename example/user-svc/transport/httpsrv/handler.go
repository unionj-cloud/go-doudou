package httpsrv

import (
	"net/http"
)

type UserHandler interface {
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

func routes(handler UserHandler) []route {
	return []route{
		{
			"SignUp",
			"POST",
			"/signup",
			handler.PostSignUp,
		},
		{
			"LogIn",
			"POST",
			"/login",
			handler.PostLogIn,
		},
		{
			"User",
			"GET",
			"/user",
			handler.GetUser,
		},
		{
			"PageUsers",
			"POST",
			"/pageusers",
			handler.PostPageUsers,
		},
	}
}
