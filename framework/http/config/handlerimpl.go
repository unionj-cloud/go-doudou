package registry

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"net/http"
	"os"
	"strings"
)

// ConfigHandlerImpl define implementation for ConfigHandler
type ConfigHandlerImpl struct {
}

// GetConfig returns all environment variables
func (receiver *ConfigHandlerImpl) GetConfig(_writer http.ResponseWriter, _req *http.Request) {
	pre := _req.FormValue("pre")
	var builder strings.Builder
	for _, pair := range os.Environ() {
		if stringutils.IsEmpty(pre) || strings.HasPrefix(pair, pre) {
			builder.WriteString(fmt.Sprintf("%s\n", pair))
		}
	}
	_writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_writer.Write([]byte(builder.String()))
}

// NewConfigHandler creates new ConfigHandlerImpl
func NewConfigHandler() ConfigHandler {
	return &ConfigHandlerImpl{}
}
