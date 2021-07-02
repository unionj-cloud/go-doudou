package constants

import "time"

const FORMAT = "2006-01-02 15:04:05"
const FORMAT2 = "2006-01-02"
const FORMAT3 = "2006/01/02"
const FORMATDOT = "2006.01.02"
const FORMATE_NANO = "2006-01-02T15:04:05.999Z"
const FORMATES = "2006-01-02T15:04:05Z"
const FORMAT4 = "2006-01-02 15:04"
const FORMAT5 = "2006年1月2日"
const FORMAT7 = "2006年01月02日"
const FORMAT6 = "2006年1月"
const FORMAT8 = "2006-01-02T15:04:05-0700" // 2020-07-12T15:31:50+0800
const FORMAT9 = "2006年01月02日15时04分"        // "2019年1月04日09时04分"
const FORMAT10 = "20060102"
const FORMAT11 = "20060102150405"
const FORMAT12 = "2006/1/2"

var (
	Loc *time.Location
)

func init() {
	var err error
	Loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}
