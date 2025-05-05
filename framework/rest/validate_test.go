package rest

import (
	"testing"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/stretchr/testify/assert"
)

// 测试结构体
type User struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}

func TestGetValidate(t *testing.T) {
	v := GetValidate()
	assert.NotNil(t, v)
	assert.IsType(t, &validator.Validate{}, v)
}

func TestGetTranslator(t *testing.T) {
	// 初始状态下translator可能为nil
	trans := GetTranslator()

	// 设置一个translator
	en := en.New()
	uni := ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	SetTranslator(trans)

	// 再次获取
	trans = GetTranslator()
	assert.NotNil(t, trans)
}

func TestSetTranslator(t *testing.T) {
	// 设置英文翻译器
	en := en.New()
	uni := ut.New(en, en)
	enTrans, _ := uni.GetTranslator("en")
	SetTranslator(enTrans)

	assert.Equal(t, enTrans, GetTranslator())

	// 设置中文翻译器
	zh := zh.New()
	uni = ut.New(zh, zh)
	zhTrans, _ := uni.GetTranslator("zh")
	SetTranslator(zhTrans)

	assert.Equal(t, zhTrans, GetTranslator())
}

func TestValidateStruct(t *testing.T) {
	// 设置英文翻译器
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	SetTranslator(trans)

	// 注册英文错误消息
	v := GetValidate()
	en_translations.RegisterDefaultTranslations(v, trans)

	// 测试有效数据
	validUser := User{
		Name:  "John",
		Email: "john@example.com",
		Age:   30,
	}
	err := ValidateStruct(validUser)
	assert.NoError(t, err)

	// 测试无效数据
	invalidUser := User{
		Name:  "",
		Email: "invalid-email",
		Age:   150,
	}
	err = ValidateStruct(invalidUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
	assert.Contains(t, err.Error(), "email")
	assert.Contains(t, err.Error(), "130")
}

func TestValidateVar(t *testing.T) {
	// 设置英文翻译器
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	SetTranslator(trans)

	// 注册英文错误消息
	v := GetValidate()
	en_translations.RegisterDefaultTranslations(v, trans)

	// 测试有效数据
	err := ValidateVar("john@example.com", "email", "")
	assert.NoError(t, err)

	// 测试无效数据
	err = ValidateVar("invalid-email", "email", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email")

	// 测试带参数的验证
	err = ValidateVar("invalid-email", "email", "Email field validation failed")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email field validation failed")
}

func TestHandleValidationErr(t *testing.T) {
	// 设置中文翻译器
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	SetTranslator(trans)

	// 注册中文错误消息
	v := GetValidate()
	zh_translations.RegisterDefaultTranslations(v, trans)

	// 生成验证错误
	user := User{}
	err := v.Struct(user)
	assert.Error(t, err)

	// 处理验证错误
	processedErr := handleValidationErr(err)
	assert.Error(t, processedErr)
	assert.Contains(t, processedErr.Error(), "必填")

	// 测试非验证错误
	nonValidationErr := assert.AnError
	processedErr = handleValidationErr(nonValidationErr)
	assert.Equal(t, nonValidationErr, processedErr)

	// 测试nil错误
	processedErr = handleValidationErr(nil)
	assert.NoError(t, processedErr)
}
