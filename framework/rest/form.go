package rest

import (
	"github.com/go-playground/form/v4"
	"net/url"
)

var decoder = form.NewDecoder()
var encoder = form.NewEncoder()

func init() {
	// frontend axios.js use [] by default
	decoder.SetNamespacePrefix("[")
	decoder.SetNamespaceSuffix("]")
	encoder.SetNamespacePrefix("[")
	encoder.SetNamespaceSuffix("]")
}

func GetFormDecoder() *form.Decoder {
	return decoder
}

func GetFormEncoder() *form.Encoder {
	return encoder
}

func DecodeForm(v interface{}, values url.Values) (err error) {
	return decoder.Decode(v, values)
}

func EncodeForm(v interface{}) (values url.Values, err error) {
	return encoder.Encode(v)
}

func RegisterFormDecoderCustomTypeFunc(fn form.DecodeCustomTypeFunc, types ...interface{}) {
	decoder.RegisterCustomTypeFunc(fn, types...)
}

func RegisterFormEncoderCustomTypeFunc(fn form.EncodeCustomTypeFunc, types ...interface{}) {
	encoder.RegisterCustomTypeFunc(fn, types...)
}
