package banner

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/unionj-cloud/go-doudou/v2/framework"
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"sync"
)

var once sync.Once

func Print() {
	once.Do(func() {
		if !framework.CheckDev() {
			return
		}
		banner := config.DefaultGddBanner
		if b, err := cast.ToBoolE(config.GddBanner.Load()); err == nil {
			banner = b
		}
		if banner {
			bannerText := config.DefaultGddBannerText
			if stringutils.IsNotEmpty(config.GddBannerText.Load()) {
				bannerText = config.GddBannerText.Load()
			}
			figure.NewColorFigure(bannerText, "doom", "green", true).Print()
		}
	})
}
