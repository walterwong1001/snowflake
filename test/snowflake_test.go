package test

import (
	"context"
	"fmt"
	"github.com/walterwong1001/snowflake/pkg/snowflake"
	clientv3 "go.etcd.io/etcd/client/v3"
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

func TestEtcd(t *testing.T) {
	config := clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}, DialTimeout: 5 * time.Second}
	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.Put(ctx, "hello", "world")
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Get(ctx, "hello")
	if err != nil {
		fmt.Println(err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s: %s", ev.Key, ev.Value)
	}
}
