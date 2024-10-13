package rest

import (
	"github.com/goccy/go-reflect"
	"github.com/unionj-cloud/toolkit/form"
	"github.com/unionj-cloud/toolkit/stringutils"
	"net/url"
)

var decoder = form.NewDecoder()
var encoder = form.NewEncoder()

func tagNameFunc(fld reflect.StructField) string {
	name := fld.Tag.Get("form")
	if stringutils.IsEmpty(name) {
		name = fld.Tag.Get("json")
	}
	return name
}

func init() {
	// frontend axios.js use [] by default
	decoder.SetNamespacePrefix("[")
	decoder.SetNamespaceSuffix("]")
	decoder.RegisterTagNameFunc(tagNameFunc)
	encoder.SetNamespacePrefix("[")
	encoder.SetNamespaceSuffix("]")
	encoder.RegisterTagNameFunc(tagNameFunc)
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
