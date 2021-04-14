package random

import (
	"fmt"
	"math/rand"
	"time"
)

// 随机数生成
// @Param	min 	int	最小值
// @Param 	max		int	最大值
// @return  string
func RandInt(min int, max int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn((max - min)) + min
	if num < min || num > max {
		RandInt(min, max)
	}
	return fmt.Sprintf("%d", num)
}
