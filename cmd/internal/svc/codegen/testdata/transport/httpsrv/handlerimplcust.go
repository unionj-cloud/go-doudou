package httpsrv

import (
	"encoding/json"
	"fmt"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"io"
	"net/http"
	"os"
)

func (receiver *UsersvcHandlerImpl) DownloadAvatar(_writer http.ResponseWriter, _req *http.Request) {
	var (
		ctx       context.Context
		userId    interface{}
		data      []byte
		userAttrs = new([]string)
		rf        *os.File
		re        error
	)
	ctx = _req.Context()
	if _req.Body == nil {
		http.Error(_writer, "missing request body", http.StatusBadRequest)
		return
	} else {
		if _err := json.NewDecoder(_req.Body).Decode(&userId); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		} else {
			if _err := ddhttp.ValidateVar(userId, "", ""); _err != nil {
				http.Error(_writer, _err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
	if _err := _req.ParseForm(); _err != nil {
		http.Error(_writer, _err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := _req.Form["data"]; exists {
		if casted, _err := cast.ToByteSliceE(_req.Form["data"]); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		} else {
			data = casted
		}
		if _err := ddhttp.ValidateVar(data, "", "data"); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if _, exists := _req.Form["data[]"]; exists {
			if casted, _err := cast.ToByteSliceE(_req.Form["data[]"]); _err != nil {
				http.Error(_writer, _err.Error(), http.StatusBadRequest)
				return
			} else {
				data = casted
			}
			if _err := ddhttp.ValidateVar(data, "", "data"); _err != nil {
				http.Error(_writer, _err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(_writer, "missing parameter data", http.StatusBadRequest)
			return
		}
	}
	if _, exists := _req.Form["userAttrs"]; exists {
		_userAttrs := _req.Form["userAttrs"]
		userAttrs = &_userAttrs
		if _err := ddhttp.ValidateVar(userAttrs, "", "userAttrs"); _err != nil {
			http.Error(_writer, _err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if _, exists := _req.Form["userAttrs[]"]; exists {
			_userAttrs := _req.Form["userAttrs[]"]
			userAttrs = &_userAttrs
			if _err := ddhttp.ValidateVar(userAttrs, "", "userAttrs"); _err != nil {
				http.Error(_writer, _err.Error(), http.StatusBadRequest)
				return
			}
		}
	}
	rf, re = receiver.usersvc.DownloadAvatar(
		ctx,
		userId,
		data,
		*userAttrs...,
	)
	if re != nil {
		if errors.Is(re, context.Canceled) {
			http.Error(_writer, re.Error(), http.StatusBadRequest)
		} else if _err, ok := re.(*ddhttp.BizError); ok {
			http.Error(_writer, _err.Error(), _err.StatusCode)
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
