package tools

import (
	"sync"
	"sync/atomic"
	"time"
)

type SnowFlake struct {
	fixedTimestamp int64 // 填充时间戳
	lastTimestamp  int64 // 上次生成的时间戳
	sequenceNumber int64 // 序列号
	machineID      int64 // 数据中心+机器号组成 共十位
	mutex          sync.Mutex
}

// NewSnowFlake machineID
func NewSnowFlake(machineID int64) *SnowFlake {
	return &SnowFlake{
		fixedTimestamp: time.Now().Unix(),
		lastTimestamp:  int64(-1),
		sequenceNumber: int64(0),
		machineID:      machineID,
	}
}

func (s *SnowFlake) GenSnowID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentTimestamp := time.Now().UnixMilli()
	if currentTimestamp == atomic.LoadInt64(&s.lastTimestamp) {
		atomic.AddInt64(&s.sequenceNumber, 1)
		if atomic.LoadInt64(&s.sequenceNumber) > 4096 {
			for currentTimestamp <= atomic.LoadInt64(&s.lastTimestamp) {
				currentTimestamp = time.Now().UnixMilli()
			}
			atomic.StoreInt64(&s.sequenceNumber, 0)
		}
	} else {
		atomic.StoreInt64(&s.sequenceNumber, 0)
	}
	atomic.StoreInt64(&s.lastTimestamp, currentTimestamp)

	return (currentTimestamp-s.fixedTimestamp)<<22 | int64(s.machineID<<12) | s.sequenceNumber
}
