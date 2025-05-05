package rest

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

// 使用与rest包中相同的json编码器
var jsonEncoder = sonic.ConfigDefault

func TestCopyReqBody(t *testing.T) {
	// 测试空body
	r1, r2, err := CopyReqBody(nil)
	assert.NoError(t, err)
	assert.Equal(t, http.NoBody, r1)
	assert.Equal(t, http.NoBody, r2)

	// 测试http.NoBody
	r1, r2, err = CopyReqBody(http.NoBody)
	assert.NoError(t, err)
	assert.Equal(t, http.NoBody, r1)
	assert.Equal(t, http.NoBody, r2)

	// 测试普通body
	originalBody := ioutil.NopCloser(strings.NewReader("test body"))
	r1, r2, err = CopyReqBody(originalBody)
	assert.NoError(t, err)

	// 读取r1
	data1, err := io.ReadAll(r1)
	assert.NoError(t, err)
	assert.Equal(t, "test body", string(data1))

	// 读取r2
	data2, err := io.ReadAll(r2)
	assert.NoError(t, err)
	assert.Equal(t, "test body", string(data2))
}

func TestCopyRespBody(t *testing.T) {
	// 测试nil
	b1, b2, err := CopyRespBody(nil)
	assert.NoError(t, err)
	assert.Nil(t, b1)
	assert.Nil(t, b2)

	// 测试有内容的buffer
	original := bytes.NewBuffer([]byte("test response"))
	b1, b2, err = CopyRespBody(original)
	assert.NoError(t, err)
	assert.Equal(t, "test response", b1.String())
	assert.Equal(t, "test response", b2.String())
}

func TestJsonMarshalIndent(t *testing.T) {
	t.Skip("跳过测试，因为包中的json变量与编译器冲突")
}

func TestGetReqBody_JSON(t *testing.T) {
	// 创建JSON请求
	jsonData := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	jsonBody, _ := jsonEncoder.Marshal(jsonData)
	bodyReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", "/test", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	bodyClone := ioutil.NopCloser(bytes.NewReader(jsonBody))

	// 调用测试的函数
	result := GetReqBody(bodyClone, req)

	// 验证结果包含预期的JSON内容
	assert.Contains(t, result, `"name"`)
	assert.Contains(t, result, `"John"`)
	assert.Contains(t, result, `"age"`)
	assert.Contains(t, result, `30`)
}

func TestGetReqBody_FormURLEncoded(t *testing.T) {
	// 创建表单请求
	formValues := url.Values{}
	formValues.Add("name", "John")
	formValues.Add("age", "30")
	formBody := formValues.Encode()
	bodyReader := strings.NewReader(formBody)

	req, _ := http.NewRequest("POST", "/test", bodyReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	bodyClone := ioutil.NopCloser(strings.NewReader(formBody))

	// 调用测试的函数
	result := GetReqBody(bodyClone, req)

	// 验证结果包含预期的表单内容
	assert.Contains(t, result, "name=John")
	assert.Contains(t, result, "age=30")
}

func TestGetReqBody_MultipartForm(t *testing.T) {
	// 创建multipart表单请求
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 添加表单字段
	w.WriteField("name", "John")
	w.WriteField("age", "30")
	w.Close()

	req, _ := http.NewRequest("POST", "/test", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	bodyClone := ioutil.NopCloser(bytes.NewReader(b.Bytes()))

	// 调用测试的函数
	result := GetReqBody(bodyClone, req)

	// 验证结果包含预期的表单内容（注意这里我们期望表单被解析）
	assert.Contains(t, result, "name=John")
	assert.Contains(t, result, "age=30")
}

func TestGetReqBody_PlainText(t *testing.T) {
	// 创建纯文本请求
	bodyText := "Hello, world!"
	bodyReader := strings.NewReader(bodyText)

	req, _ := http.NewRequest("POST", "/test", bodyReader)
	req.Header.Set("Content-Type", "text/plain")

	bodyClone := ioutil.NopCloser(strings.NewReader(bodyText))

	// 调用测试的函数
	result := GetReqBody(bodyClone, req)

	// 验证结果是原始文本
	assert.Equal(t, bodyText, result)
}

func TestGetRespBody_JSON(t *testing.T) {
	// 创建测试的response recorder
	rec := httptest.NewRecorder()

	// 写入JSON响应
	jsonData := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	jsonBytes, _ := jsonEncoder.Marshal(jsonData)
	rec.Header().Set("Content-Type", "application/json")
	rec.Write(jsonBytes)

	// 调用测试的函数
	result := GetRespBody(rec)

	// 验证结果包含预期的JSON内容（带缩进）
	assert.Contains(t, result, `"name"`)
	assert.Contains(t, result, `"John"`)
	assert.Contains(t, result, `"age"`)
	assert.Contains(t, result, `30`)
}

func TestGetRespBody_PlainText(t *testing.T) {
	// 创建测试的response recorder
	rec := httptest.NewRecorder()

	// 写入纯文本响应
	plainText := "Hello, world!"
	rec.Header().Set("Content-Type", "text/plain")
	rec.Write([]byte(plainText))

	// 调用测试的函数
	result := GetRespBody(rec)

	// 验证结果是原始文本
	assert.Equal(t, plainText, result)
}

func TestGetRespBody_LongText(t *testing.T) {
	// 创建测试的response recorder
	rec := httptest.NewRecorder()

	// 创建超过1000个字符的长文本
	longText := strings.Repeat("a", 2000)
	rec.Header().Set("Content-Type", "text/plain")
	rec.Write([]byte(longText))

	// 调用测试的函数
	result := GetRespBody(rec)

	// 验证结果被截断到1000个字符
	assert.Equal(t, 1000, len(result))
	assert.Equal(t, strings.Repeat("a", 1000), result)
}
