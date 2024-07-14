package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	WORK_ID_BITS  = 10
	SEQUENCE_BITS = 12
	MAX_WORK_ID   = -1 ^ (-1 << WORK_ID_BITS)
	SEQUENCE_MASK = -1 ^ (-1 << SEQUENCE_BITS)
	EPOCH         = 1704067200000 // 2024-01-01 00:00:00 UTC in milliseconds
)

type snowflake struct {
	mutex    sync.Mutex
	epoch    int64
	workId   uint16 // 占用10位
	sequence uint64 // 占用12位
	elapsed  int64
}

func NewSnowflake(workId uint16) (*snowflake, error) {
	if workId > MAX_WORK_ID {
		return nil, errors.New("workId exceeds max value")
	}
	return &snowflake{
		epoch:    EPOCH,
		workId:   workId,
		sequence: 0,
		elapsed:  0,
	}, nil
}

func (s *snowflake) NextID() uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	currentMilli := time.Now().UnixNano() / 1000000 // 转换为毫秒

	elapsedMilli := currentMilli - s.epoch
	if s.elapsed < elapsedMilli {
		s.elapsed = elapsedMilli
		s.sequence = 0
	} else {
		s.sequence = (s.sequence + 1) & SEQUENCE_MASK
		if s.sequence == 0 {
			s.elapsed++
			// 等待下一毫秒
			for elapsedMilli <= s.elapsed {
				elapsedMilli = (time.Now().UnixNano() / 1000000) - s.epoch
			}
		}
	}

	id := uint64(currentMilli-s.epoch)<<(WORK_ID_BITS+SEQUENCE_BITS) | uint64(s.workId)<<SEQUENCE_BITS | s.sequence

	return id
}
