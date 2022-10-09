package model

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	logger "github.com/unionj-cloud/go-doudou/toolkit/zlogger"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type Annotation struct {
	Name   string
	Params []string
}

type AnnotationStore map[string][]Annotation

func (receiver AnnotationStore) HasAnnotation(key string, annotationName string) bool {
	for _, item := range receiver[key] {
		if item.Name == annotationName {
			return true
		}
	}
	return false
}

func (receiver AnnotationStore) GetParams(key string, annotationName string) []string {
	for _, item := range receiver[key] {
		if item.Name == annotationName {
			return item.Params
		}
	}
	return nil
}

// Route wraps config for route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// borrowed from httputil unexported function drainBody
func CopyReqBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func CopyRespBody(b *bytes.Buffer) (b1, b2 *bytes.Buffer, err error) {
	if b == nil {
		return
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	return &buf, bytes.NewBuffer(buf.Bytes()), nil
}

func JsonMarshalIndent(data interface{}, prefix, indent string, disableHTMLEscape bool) (string, error) {
	b := &bytes.Buffer{}
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!disableHTMLEscape)
	encoder.SetIndent(prefix, indent)
	if err := encoder.Encode(data); err != nil {
		return "", errors.Errorf("failed to marshal data to JSON, %s", err)
	}
	return b.String(), nil
}

func GetReqBody(cp io.ReadCloser, r *http.Request) string {
	var contentType string
	if len(r.Header["Content-Type"]) > 0 {
		contentType = r.Header["Content-Type"][0]
	}
	var reqBody string
	if cp != nil {
		if strings.Contains(contentType, "multipart/form-data") {
			r.Body = cp
			if err := r.ParseMultipartForm(32 << 20); err == nil {
				reqBody = r.Form.Encode()
				if unescape, err := url.QueryUnescape(reqBody); err == nil {
					reqBody = unescape
				}
			} else {
				logger.Error().Err(err).Msg("call r.ParseMultipartForm(32 << 20) error")
			}
		} else if strings.Contains(contentType, "application/json") {
			data := make(map[string]interface{})
			if err := json.NewDecoder(cp).Decode(&data); err == nil {
				b, _ := json.MarshalIndent(data, "", "    ")
				reqBody = string(b)
			} else {
				logger.Error().Err(err).Msg("call json.NewDecoder(reqBodyCopy).Decode(&data) error")
			}
		} else {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(cp); err == nil {
				data := []rune(buf.String())
				end := len(data)
				if end > 1000 {
					end = 1000
				}
				reqBody = string(data[:end])
				if strings.Contains(contentType, "application/x-www-form-urlencoded") {
					if unescape, err := url.QueryUnescape(reqBody); err == nil {
						reqBody = unescape
					}
				}
			} else {
				logger.Error().Err(err).Msg("call buf.ReadFrom(reqBodyCopy) error")
			}
		}
	}
	return reqBody
}

func GetRespBody(rec *httptest.ResponseRecorder) string {
	var (
		respBody string
		err      error
	)
	if strings.Contains(rec.Result().Header.Get("Content-Type"), "application/json") {
		var respBodyCopy *bytes.Buffer
		if respBodyCopy, rec.Body, err = CopyRespBody(rec.Body); err == nil {
			data := make(map[string]interface{})
			if err := json.NewDecoder(rec.Body).Decode(&data); err == nil {
				b, _ := json.MarshalIndent(data, "", "    ")
				respBody = string(b)
			} else {
				logger.Error().Err(err).Msg("call json.NewDecoder(rec.Body).Decode(&data) error")
			}
		} else {
			logger.Error().Err(err).Msg("call respBodyCopy.ReadFrom(rec.Body) error")
		}
		rec.Body = respBodyCopy
	} else {
		data := []rune(rec.Body.String())
		end := len(data)
		if end > 1000 {
			end = 1000
		}
		respBody = string(data[:end])
	}
	return respBody
}
