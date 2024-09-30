package snowflake

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// time - 41 bits (millisecond precision w/ a custom epoch gives us 69 years)
// configured machine id - 10 bits - gives us up to 1024 machines
// sequence number - 12 bits - rolls over every 4096 per machine (with protection to avoid rollover in the same ms)
const (
	MACHINE_ID_BITS = 10                           // bit length of machine id
	SEQUENCE_BITS   = 12                           // bit length of sequence number
	MAX_MACHINE_ID  = -1 ^ (-1 << MACHINE_ID_BITS) // max machine id (1023)
	SEQUENCE_MASK   = -1 ^ (-1 << SEQUENCE_BITS)   // sequence mask 4095
	EPOCH           = 1721001600000                // 2024-07-15 00:00:00 UTC in milliseconds
)

type Snowflake struct {
	mutex         sync.Mutex
	machineId     int
	sequence      int
	lastTimestamp int64
}

func NewSnowflake(machineId int) (*Snowflake, error) {
	if machineId > MAX_MACHINE_ID {
		return nil, errors.New(fmt.Sprintf("Machine Id can't be greater than %d", MAX_MACHINE_ID))
	}
	return &Snowflake{
		machineId:     machineId,
		sequence:      0,
		lastTimestamp: 0,
	}, nil
}

func (s *Snowflake) NextID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixMilli() // 转换为毫秒
	if timestamp < s.lastTimestamp {
		log.Printf("clock is moving backwards.  Rejecting requests until %d.\n", s.lastTimestamp)
		return 0
	}

	if s.lastTimestamp == timestamp {
		s.sequence = (s.sequence + 1) & SEQUENCE_MASK
		if s.sequence == 0 {
			timestamp = tilNextMillis(timestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return (timestamp-EPOCH)<<(MACHINE_ID_BITS+SEQUENCE_BITS) | int64(s.machineId<<SEQUENCE_BITS) | int64(s.sequence)
}

func tilNextMillis(timestamp int64) int64 {
	var current = time.Now().UnixMilli()
	for current <= timestamp {
		current = time.Now().UnixMilli()
	}
	return current
}
