package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
)

type Lock struct {
	Client  *clientv3.Client
	mutex   *concurrency.Mutex
	session *concurrency.Session
	mu      sync.Mutex
}

// Lock 获取锁
func (l *Lock) Lock(ctx context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.session == nil {
		// 创建一个 session，并保存在结构体中
		session, err := concurrency.NewSession(l.Client)
		if err != nil {
			return err
		}
		l.session = session
	}

	// 创建锁对象并尝试加锁
	l.mutex = concurrency.NewMutex(l.session, key)
	return l.mutex.Lock(ctx)
}

// Unlock 释放锁
func (l *Lock) Unlock(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.mutex == nil {
		return fmt.Errorf("mutex is not initialized")
	}
	return l.mutex.Unlock(ctx)
}

// Close 关闭 Session，释放资源
func (l *Lock) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.session != nil {
		return l.session.Close()
	}
	return nil
}
