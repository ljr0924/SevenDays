package cache

import (
    "SevenDays/cache/container"
    "sync"
)

// 添加互斥锁mu，支持并发
type cache struct {
    mu         sync.Mutex
    c          container.Container
    cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
    c.mu.Lock()
    defer c.mu.Unlock()
    // 延迟初始化，提高程序启动速度
    if c.c == nil {
        // todo: 改成工厂函数获取缓存器
        c.c = container.NewLru(c.cacheBytes, nil)
    }
    c.c.Add(key, value)
}

func (c *cache) get(key string) (bv ByteView, ok bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.c == nil {
        return
    }
    if v, ok := c.c.Get(key);ok {
        return v.(ByteView), ok
    }
    return
}