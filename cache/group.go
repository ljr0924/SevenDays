package cache

import (
    "fmt"
    "log"
    "sync"
)

// 缓存组，实现类似命名空间功能，每一个组就是一个命名空间
type Group struct {
    name      string
    getter    Getter
    mainCache cache
}

var (
    mu sync.RWMutex
    groups = make(map[string]*Group)
)

// Group构造器
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
    if getter == nil {
        panic("getter is required")
    }
    mu.Lock()
    defer mu.Unlock()
    g := &Group{
        name:      name,
        getter:    getter,
        mainCache: cache{cacheBytes: cacheBytes},
    }
    groups[name] = g

    return g
}

// 获取缓存组，不存在返回nil
func GetGroup(name string) *Group {
    mu.RLock()
    g:= groups[name]
    mu.RUnlock()
    return g
}

func (g *Group) Get(key string) (ByteView, error) {
    if key == "" {
        return ByteView{}, fmt.Errorf("key is required")
    }

    if v, ok := g.mainCache.get(key); ok {
        log.Println("hit cache")
        return v, nil
    }

    return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
    return g.getLocally(key)
}

// 调用用户自定义回调函数获取数据，并更新到cache中
// 有效防止缓存穿透
func (g *Group) getLocally(key string) (value ByteView, err error) {
    bytes, err := g.getter.Get(key)
    if err != nil {
        return
    }

    value = ByteView{
        b: cloneBytes(bytes),
    }
    g.populateCache(key, value)
    return value, nil
}

// 更新缓存
func (g *Group) populateCache(key string, value ByteView) {
    g.mainCache.add(key, value)
}
