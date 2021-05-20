package httpsrv

import (
	"net/http"
	service "usersvc"
)

type UserServiceHandlerImpl struct {
	userService service.UserService
}

func (receiver *UserServiceHandlerImpl) PostSignUp(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
func (receiver *UserServiceHandlerImpl) PostLogIn(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
func (receiver *UserServiceHandlerImpl) GetUser(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
func (receiver *UserServiceHandlerImpl) PostPageUsers(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func NewUserServiceHandler(userService service.UserService) UserServiceHandler {
	return &UserServiceHandlerImpl{
		userService,
	}
}
