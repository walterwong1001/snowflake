package test

import (
	"context"
	"github.com/walterwong1001/snowflake/pkg/snowflake"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
	"time"
)

func TestIdGenerator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	config := clientv3.Config{Endpoints: []string{"localhost:2179", "localhost:2279", "localhost:2379"}, DialTimeout: 5 * time.Second}
	for i := 0; i < 2; i++ {
		go func() {
			gen := snowflake.NewGenerator(ctx, config)
			log.Println(gen.Next())
		}()
	}
	time.Sleep(2 * time.Minute)
	cancel()
	log.Println("======================================")
}
