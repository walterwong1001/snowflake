package snowflake

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/walterwong1001/snowflake/pkg/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	MAX_WORK_ID        = (1 << 10) - 1 // 1023 = 2^10 - 1
	WORK_ID_KEY_PREFIX = "/snowflake/worker-"
	LOCK_KEY           = "/snowflake/lock"
)

// 注册并返回Work id
func workId(ctx context.Context, client *clientv3.Client) (int, error) {
	seq, err := register(ctx, client)
	if err != nil {
		return -1, err
	}
	return seq, nil
}

// 获取ETCD中已经存在的Work id
func getExistingWorkerIDs(ctx context.Context, client *clientv3.Client) (map[int]int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, WORK_ID_KEY_PREFIX, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	m := make(map[int]int)
	for _, e := range resp.Kvs {
		seq, _ := strconv.Atoi(string(e.Value))
		m[seq] = seq
	}
	return m, nil
}

// 找到可用的Work id
func findAvailableWorkerID(m map[int]int) (int, error) {
	for i := 0; i <= MAX_WORK_ID; i++ {
		if _, ok := m[i]; !ok {
			return i, nil
		}
	}
	return -1, errors.New("no available work id")
}

// 注册Work id
func register(ctx context.Context, client *clientv3.Client) (int, error) {
	lock := &etcd.Lock{Client: client}
	defer lock.Close()

	if err := lock.Lock(ctx, LOCK_KEY); err != nil {
		return -1, err
	}
	defer func() {
		if err := lock.Unlock(ctx); err != nil {
			log.Printf("Failed to unlock: %v", err)
		}
	}()

	return registerWorkIDWithLease(ctx, client)
}

// 注册租约
func registerWorkIDWithLease(ctx context.Context, client *clientv3.Client) (int, error) {
	m, err := getExistingWorkerIDs(ctx, client)
	if err != nil {
		return -1, err
	}
	id, err := findAvailableWorkerID(m)
	if err != nil {
		return -1, err
	}

	resp, err := clientv3.NewLease(client).Grant(ctx, 60)
	if err != nil {
		return -1, err
	}
	v := strconv.Itoa(id)
	_, err = client.Put(ctx, WORK_ID_KEY_PREFIX+v, v, clientv3.WithLease(resp.ID))
	if err != nil {
		return -1, err
	}

	go schedule(ctx, client, WORK_ID_KEY_PREFIX+v, v)

	return id, nil
}

// 定时续约
func schedule(ctx context.Context, client *clientv3.Client, key, value string) {

	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stop to renew lease: %s : %s", key, value)
			ticker.Stop()
			return
		case <-ticker.C:
			if err := renewLease(ctx, client, key, value); err != nil {
				log.Printf("Failed to renew lease: %s : %s", key, value)
			}
		}
	}
}

// 续约
func renewLease(ctx context.Context, client *clientv3.Client, key, value string) error {
	resp, err := clientv3.NewLease(client).Grant(ctx, 60)
	if err != nil {
		return err
	}

	_, err = client.Put(ctx, key, value, clientv3.WithLease(resp.ID))

	return err
}
