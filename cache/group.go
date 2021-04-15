package cache

import (
    pb "SevenDays/cache/cache_pb"
    "SevenDays/cache/single_flight"
    "fmt"
    "log"
    "sync"
)

// 缓存组，实现类似命名空间功能，每一个组就是一个命名空间
type Group struct {
    name      string
    getter    Getter
    mainCache cache
    peer      PeerPicker

    loader *single_flight.Group
}

var (
    mu     sync.RWMutex
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
        loader:    &single_flight.Group{},
    }
    groups[name] = g

    return g
}

// 获取缓存组，不存在返回nil
func GetGroup(name string) *Group {
    mu.RLock()
    g := groups[name]
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

    // 使用 g.loader.Do 包裹起来即可，这样确保了并发场景下针对相同的 key，load 过程只会调用一次
    view, err := g.loader.Do(key, func() (interface{}, error) {
        // 若非本机节点，则调用getFromPeer()从远程获取。
        if g.peer != nil {
            if peer, ok := g.peer.PickPeer(key); ok {
                if value, err = g.getFromPeer(peer, key); err == nil {
                    return value, nil
                }
                log.Println("[Cache] Failed to get from peer", err)
            }
        }
        // 若是本机节点或者远程获取失败，调用getLocally
        return g.getLocally(key)
    })

    if err == nil {
        return view.(ByteView), nil
    }

    return ByteView{}, err
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

// 将实现PeerPicker接口的HTTPPool注入到Group中
func (g *Group) RegisterPeers(peer PeerPicker) {
    if g.peer != nil {
        panic("RegisterPeerPicker called more than once")
    }
    g.peer = peer
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
    req := &pb.Request{
        Group: g.name,
        Key:   key,
    }
    res := &pb.Response{}
    err := peer.Get(req, res)
    if err != nil {
        return ByteView{}, err
    }
    return ByteView{b: res.Value}, nil
}
