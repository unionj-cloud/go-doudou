package main

//import (
//	"bytes"
//	"context"
//	"github.com/goccy/go-json"
//	"fmt"
//	"github.com/go-playground/form/v4"
//	"github.com/unionj-cloud/go-doudou/toolkit/lab/vo"
//	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
//	"net/url"
//	"os"
//	"reflect"
//)
//
//type Usersvc interface {
//	PageUsers(ctx context.Context, query *vo.PageQuery) (data vo.PageRet, err error)
//
//	GetUser(ctx context.Context, userId int) (data vo.UserVo, err error)
//
//	PublicSignUp(ctx context.Context, username string, password string, code *string) (data int, err error)
//
//	PublicLogIn(ctx context.Context, username string, password string) (data string, err error)
//
//	UploadAvatar(ctx context.Context, avatar v3.FileModel, id int) (data string, err error)
//
//	GetPublicDownloadAvatar(ctx context.Context, userId int) (data *os.File, err error)
//}
//
//type UsersvcImpl struct {
//}
//
//func (u *UsersvcImpl) PageUsers(ctx context.Context, query *vo.PageQuery) (data vo.PageRet, err error) {
//	data.PageSize = query.Page.Size
//	//err = errors.New("special error")
//	return
//}
//
//func (u *UsersvcImpl) GetUser(ctx context.Context, userId int) (data vo.UserVo, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (u *UsersvcImpl) PublicSignUp(ctx context.Context, username string, password string, code *string) (data int, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (u *UsersvcImpl) PublicLogIn(ctx context.Context, username string, password string) (data string, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (u *UsersvcImpl) UploadAvatar(ctx context.Context, avatar v3.FileModel, id int) (data string, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (u *UsersvcImpl) GetPublicDownloadAvatar(ctx context.Context, userId int) (data *os.File, err error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func main() {
//	svc := &UsersvcImpl{}
//	reflectType := reflect.TypeOf(svc)
//	for i := 0; i < reflectType.NumMethod(); i++ {
//		m := reflectType.Method(i)
//		fmt.Print(m.Name)
//		fmt.Print("\t")
//		for j := 0; j < m.Type.NumIn(); j++ {
//			pType := m.Type.In(j)
//			fmt.Print(pType.String())
//			fmt.Print("\t")
//		}
//		fmt.Println(m.Type.String())
//	}
//
//	pqPtr := reflect.New(TypeRegistry["vo.PageQuery"])
//	expect := vo.PageQuery{
//		Page: vo.Page{
//			Orders: nil,
//			PageNo: 1,
//			Size:   10,
//		},
//	}
//	j, _ := json.Marshal(expect)
//	_ = json.NewDecoder(bytes.NewReader(j)).Decode(pqPtr.Interface())
//	med, ok := reflectType.MethodByName("PageUsers")
//	if ok {
//		bodyType := med.Type.In(2)
//		fmt.Println(bodyType.String())
//		fmt.Println(bodyType.Kind() == reflect.Ptr)
//		bodyValue := pqPtr.Elem()
//		if bodyType.Kind() == reflect.Ptr {
//			bodyValue = pqPtr
//		}
//		resultValues := med.Func.Call([]reflect.Value{reflect.ValueOf(svc), reflect.ValueOf(context.Background()), bodyValue})
//		jsonRet, _ := json.Marshal(resultValues[0].Interface())
//		fmt.Println(string(jsonRet))
//		err := resultValues[1].Interface()
//		fmt.Println(err)
//	}
//
//	u := vo.UserVo{
//		Username: "wubin1989",
//		Name:     "武斌",
//	}
//	j, _ = json.Marshal(u)
//
//	queryParams := make(url.Values)
//	queryParams.Set("id", "12")
//	queryParams.Set("phone", "13552053960")
//	queryParams.Set("dept", "tech")
//	queryParams.Set("name", "斌斌")
//
//	decoder := form.NewDecoder()
//	var copyUser vo.UserVo
//	err := decoder.Decode(&copyUser, queryParams)
//	if err != nil {
//		panic(err)
//	}
//	err = json.Unmarshal(j, &copyUser)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(copyUser)
//}
