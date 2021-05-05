package httpsrv

import (
	service "example/user-svc"
	"github.com/gorilla/schema"
	"net/http"
)

type userHandlerImpl struct {
	decoder     *schema.Decoder
	userService service.UserService
}

func (u userHandlerImpl) PostSignUp(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (u userHandlerImpl) PostLogIn(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (u userHandlerImpl) GetUser(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (u userHandlerImpl) PostPageUsers(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewUserHandler(usersvc service.UserService) UserHandler {
	return userHandlerImpl{
		decoder:     schema.NewDecoder(),
		userService: usersvc,
	}
}
