package service

import (
	"GoLeaf/config"
	"sync/atomic"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
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

func (sf *Snowflake) GenerateID() int64 {
	now := time.Now().UnixNano() / 1e6
	//todo 时钟回拨处理
	//if now < sf.latestTime {
	//	panic("Clock moved backwards, refusing to generate id")
	//}
	// 如果当前时间与上次生成ID的时间相同，则进行毫秒内序列
	if now == atomic.LoadInt64(&sf.latestTime) {
		for {
			sequenceID := atomic.LoadInt64(&sf.sequenceID)
			nextSequenceID := (sequenceID + 1) & 4095
			// 如果当前毫秒内序列已经用完，则等待下一毫秒
			if nextSequenceID == 0 {
				now = time.Now().UnixNano() / 1e6
				for now <= sf.latestTime {
					now = time.Now().UnixNano() / 1e6
				}
			}
			if atomic.CompareAndSwapInt64(&sf.sequenceID, sequenceID, nextSequenceID) {
				break
			}
		}
	} else {
		atomic.StoreInt64(&sf.sequenceID, 0)
	}
	atomic.StoreInt64(&sf.latestTime, now)
	ID := (now-config.StartTime)<<22 | (sf.workerID << 12) | atomic.LoadInt64(&sf.sequenceID)
	return ID
}
