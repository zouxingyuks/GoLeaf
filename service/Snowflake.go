package service

import (
	"GoLeaf/config"
	"sync"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
	mu         sync.Mutex
	latestTime int64
	workerID   int64
	sequenceID int64
}

// NewSnowflake 创建一个新的 Snowflake 实例
func NewSnowflake(workerID int64) *Snowflake {
	return &Snowflake{
		latestTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6,
		workerID:   workerID,
		sequenceID: 0,
	}
}

// GenerateID 生成唯一标识符
func (sf *Snowflake) GenerateID() int64 {

	now := time.Now().UnixNano() / 1e6
	//todo 时钟回拨处理
	//if now < sf.latestTime {
	//	panic("Clock moved backwards, refusing to generate id")
	//}
	// 如果当前时间与上次生成ID的时间相同，则进行毫秒内序列
	if now == sf.latestTime {
		sf.mu.Lock()

		sf.sequenceID = (sf.sequenceID + 1) & 4095
		if sf.sequenceID == 0 {
			for now <= sf.latestTime {
				now = time.Now().UnixNano() / 1e6
			}
		}
		sf.mu.Unlock()
	} else {
		sf.sequenceID = 0
	}

	sf.latestTime = now
	ID := (now-config.StartTime)<<22 | (sf.workerID << 12) | sf.sequenceID
	return ID
}
