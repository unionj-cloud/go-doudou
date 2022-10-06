package banner

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/unionj-cloud/go-doudou/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/toolkit/stringutils"
	"sync"
)

var once sync.Once

func Print() {
	once.Do(func() {
		if !config.CheckDev() {
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
