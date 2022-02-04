package httpsrv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/cast"
	v3 "github.com/unionj-cloud/go-doudou/openapi/v3"
	"io"
	"mime/multipart"
	"os"
	service "testdata"
	"net/http"
	"testdata/vo"
	"github.com/pkg/errors"
)

type UsersvcHandlerImpl struct{
	usersvc service.Usersvc
}
func (receiver *UsersvcHandlerImpl) PageUsers(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx context.Context
		query vo.PageQuery
		code int
		data vo.PageRet
		msg error
	)
	ctx = _req.Context()
	if _req.Body == nil {
		http.Error(_writer, "missing request body", http.StatusBadRequest)
		return
	} else {
		if err := json.NewDecoder(_req.Body).Decode(&query); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
	code,data,msg = receiver.usersvc.PageUsers(
		ctx,
		query,
	)
	if msg != nil {
		if errors.Is(msg, context.Canceled) {
			http.Error(_writer, msg.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, msg.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(_writer).Encode(struct{
		Code int `json:"code,omitempty"`
		Data vo.PageRet `json:"data,omitempty"`
	}{
		Code: code,
		Data: data,
	}); err != nil {
		http.Error(_writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (receiver *UsersvcHandlerImpl) GetUser(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx context.Context
		userId string
		photo string
		code int
		data string
		msg error
	)
	ctx = _req.Context()
	if err := _req.ParseForm(); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["userId"]; exists {
		userId = _req.FormValue("userId")
	} else {
		http.Error(_writer, "missing parameter userId", http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["photo"]; exists {
		photo = _req.FormValue("photo")
	} else {
		http.Error(_writer, "missing parameter photo", http.StatusBadRequest)
		return
	}
	code,data,msg = receiver.usersvc.GetUser(
		ctx,
		userId,
		photo,
	)
	if msg != nil {
		if errors.Is(msg, context.Canceled) {
			http.Error(_writer, msg.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, msg.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(_writer).Encode(struct{
		Code int `json:"code,omitempty"`
		Data string `json:"data,omitempty"`
	}{
		Code: code,
		Data: data,
	}); err != nil {
		http.Error(_writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (receiver *UsersvcHandlerImpl) SignUp(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx context.Context
		username string
		password int
		actived bool
		score []int
		code int
		data string
		msg error
	)
	ctx = _req.Context()
	if err := _req.ParseForm(); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["username"]; exists {
		username = _req.FormValue("username")
	} else {
		http.Error(_writer, "missing parameter username", http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["password"]; exists {
		if casted, err := cast.ToIntE(_req.FormValue("password")); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		} else {
			password = casted
		}
	} else {
		http.Error(_writer, "missing parameter password", http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["actived"]; exists {
		if casted, err := cast.ToBoolE(_req.FormValue("actived")); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		} else {
			actived = casted
		}
	} else {
		http.Error(_writer, "missing parameter actived", http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["score"]; exists {
		if casted, err := cast.ToIntSliceE(_req.Form["score"]); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		} else {
			score = casted
		}
	} else {
		if _, exists := _req.Form["score[]"]; exists {
			if casted, err := cast.ToIntSliceE(_req.Form["score[]"]); err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			} else {
				score = casted
			}
		} else {
			http.Error(_writer, "missing parameter score", http.StatusBadRequest)
			return
		}
	}
	code,data,msg = receiver.usersvc.SignUp(
		ctx,
		username,
		password,
		actived,
		score,
	)
	if msg != nil {
		if errors.Is(msg, context.Canceled) {
			http.Error(_writer, msg.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, msg.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(_writer).Encode(struct{
		Code int `json:"code,omitempty"`
		Data string `json:"data,omitempty"`
	}{
		Code: code,
		Data: data,
	}); err != nil {
		http.Error(_writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (receiver *UsersvcHandlerImpl) UploadAvatar(_writer http.ResponseWriter, _req *http.Request) {
	var (
		pc context.Context
		pf []v3.FileModel
		ps string
		pf2 v3.FileModel
		pf3 *multipart.FileHeader
		pf4 []*multipart.FileHeader
		ri int
		rs string
		re error
	)
	pc = _req.Context()
	if err := _req.ParseMultipartForm(32 << 20); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	pfFileHeaders, exists := _req.MultipartForm.File["pf"]
	if exists {
		if len(pfFileHeaders) == 0 {
			http.Error(_writer, "no file uploaded for parameter pf", http.StatusBadRequest)
			return
		}
		for _, _fh :=range pfFileHeaders {
			_f, err := _fh.Open()
			if err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			}
			pf = append(pf, v3.FileModel{
				Filename: _fh.Filename,
				Reader: _f,
			})
		}
	} else {
		http.Error(_writer, "missing parameter pf", http.StatusBadRequest)
		return
	}
	if err := _req.ParseForm(); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["ps"]; exists {
		ps = _req.FormValue("ps")
	} else {
		http.Error(_writer, "missing parameter ps", http.StatusBadRequest)
		return
	}
	pf2FileHeaders, exists := _req.MultipartForm.File["pf2"]
	if exists {
		if len(pf2FileHeaders) == 0 {
			http.Error(_writer, "no file uploaded for parameter pf2", http.StatusBadRequest)
			return
		}
		if len(pf2FileHeaders) > 0 {
			_fh := pf2FileHeaders[0]
			_f, err := _fh.Open()
			if err != nil {
				http.Error(_writer, err.Error(), http.StatusBadRequest)
				return
			}
			pf2 = v3.FileModel{
				Filename: _fh.Filename,
				Reader: _f,
			}
		}
	} else {
		http.Error(_writer, "missing parameter pf2", http.StatusBadRequest)
		return
	}
	pf3Files := _req.MultipartForm.File["pf3"]
	if len(pf3Files) > 0 {
		pf3 = pf3Files[0]
	}
	pf4 = _req.MultipartForm.File["pf4"]
	ri,rs,re = receiver.usersvc.UploadAvatar(
		pc,
		pf,
		ps,
		pf2,
		pf3,
		pf4,
	)
	if re != nil {
		if errors.Is(re, context.Canceled) {
			http.Error(_writer, re.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, re.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(_writer).Encode(struct{
		Ri int `json:"ri,omitempty"`
		Rs string `json:"rs,omitempty"`
	}{
		Ri: ri,
		Rs: rs,
	}); err != nil {
		http.Error(_writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (receiver *UsersvcHandlerImpl) DownloadAvatar(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx context.Context
		userId string
		userAttrs []string
	rf *os.File
	re error
	)
	ctx = _req.Context()
	if err := _req.ParseForm(); err != nil {
		http.Error(_writer, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["userId"]; exists {
		userId = _req.FormValue("userId")
	} else {
		http.Error(_writer, "missing parameter userId", http.StatusBadRequest)
		return
	}
	if _req.Body == nil {
		http.Error(_writer, "missing request body", http.StatusBadRequest)
		return
	} else {
		if err := json.NewDecoder(_req.Body).Decode(&userAttrs); err != nil {
			http.Error(_writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
	rf,re = receiver.usersvc.DownloadAvatar(
		ctx,
		userId,
		userAttrs...,
	)
	if re != nil {
		if errors.Is(re, context.Canceled) {
			http.Error(_writer, re.Error(), http.StatusBadRequest)
		} else {
			http.Error(_writer, re.Error(), http.StatusInternalServerError)
		}
		return
	}
	if rf == nil {
		http.Error(_writer, "No file returned", http.StatusInternalServerError)
		return
	}
	defer rf.Close()
	var _fi os.FileInfo
	_fi, _err := rf.Stat()
	if _err != nil {
		http.Error(_writer, _err.Error(), http.StatusInternalServerError)
		return
	}
	_writer.Header().Set("Content-Disposition", "attachment; filename="+_fi.Name())
	_writer.Header().Set("Content-Type", "application/octet-stream")
	_writer.Header().Set("Content-Length", fmt.Sprintf("%d", _fi.Size()))
	io.Copy(_writer, rf)
}


func NewUsersvcHandler(usersvc service.Usersvc) UsersvcHandler {
	return &UsersvcHandlerImpl{
		usersvc,
	}
}