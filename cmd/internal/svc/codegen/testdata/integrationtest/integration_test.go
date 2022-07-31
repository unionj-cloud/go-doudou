package integrationtest

import (
	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
	"net/http"
	"testing"
)

var (
	router *mux.Router
)

func TestMain(m *testing.M) {
	m.Run()
}

func Test_测试1(t *testing.T) {
	apitest.New().
		Handler(router).
		Get("/user/1234").
		Query("userId", "2621936").
		Expect(t).
		Body(`{"id": "1234", "name": "Tate"}`).
		Status(http.StatusOK).
		End()
}

func Test_测试2(t *testing.T) {
	apitest.New().
		Handler(router).
		Get("/user/1234").
		Query("userId", "2863467").
		Expect(t).
		Body(`{"id": "4321", "name": "Tate"}`).
		Status(http.StatusOK).
		End()
}
