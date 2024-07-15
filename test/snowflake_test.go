package test

import (
	"github.com/walterwong1001/snowflake/pkg/snowflake"
	"log"
	"testing"
	"time"
)

func TestNewSnowflake(t *testing.T) {
	s, err := snowflake.NewSnowflake(1)
	if err != nil {
		log.Println(err)
		return
	}
	var count int64

	// 获取当前时间
	startTime := time.Now()

	// 在一秒内尽可能多地调用 NextID() 并计数
	for {
		// 检查是否超过一秒
		if time.Since(startTime) >= time.Second {
			break
		}
		s.NextID()

		count++
	}

	// 输出调用次数
	log.Println("Count:", count)
}
