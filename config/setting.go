package config

import (
	"time"
)

var (
	StartTime int64 = 1577808000000 // 默认值为2020-01-01
)

func loadSetting() {
	now, err := time.Parse("2006-01-02", Configs.GetString("goLeaf.snowflake.start"))
	if err != nil {
		//todo 时间戳设置错误处理
		panic(err)
	}
	StartTime = now.UnixNano() / 1e6
}
