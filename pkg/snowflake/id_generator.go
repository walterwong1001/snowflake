package snowflake

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type IdGenerator interface {
	Next() int64
	WorkId() int
}

type etcdIdGenerator struct {
	etcdConfig clientv3.Config
	workId     int
	snowflake  *Snowflake
}

func NewGenerator(ctx context.Context, conf clientv3.Config) IdGenerator {
	generator := &etcdIdGenerator{etcdConfig: conf}
	client, err := clientv3.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	workId, err := workId(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	generator.workId = workId
	snowflake, err := NewSnowflake(workId)
	if err != nil {
		log.Fatal(err)
	}
	generator.snowflake = snowflake

	return generator
}

func (g *etcdIdGenerator) Next() int64 {
	return g.snowflake.NextID()
}

func (g *etcdIdGenerator) WorkId() int {
	return g.workId
}
