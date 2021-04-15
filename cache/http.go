package cache

import (
    pb "SevenDays/cache/cache_pb"
    "google.golang.org/protobuf/proto"
    "log"
    "net/http"
    "strings"
    "sync"
)

const defaultBasePath = "/_banana"
const defaultReplicas = 50

type HTTPPool struct {
    self        string // 用来记录自己的地址，包括主机名/IP 和端口。
    basePath    string
    mu          sync.Mutex
    peers       *Map
    /*
    每一个远程节点对应一个httpGetter
    keyed by e.g. "http://10.0.0.2:8008"
    */
    httpGetters map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
    return &HTTPPool{
        self:     self,
        basePath: defaultBasePath,
    }
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
    log.Printf(format, v...)
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 1. 判断前缀
    if !strings.HasPrefix(r.URL.Path, p.basePath) {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    p.Log("%s %s", r.Method, r.URL.Path)

    // 2. 切割url  /base url/{group}/{key}
    parts := strings.SplitN(r.URL.Path[len(p.basePath)+1:], "/", 2)
    if len(parts) != 2 {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    groupName := parts[0]
    key := parts[1]

    p.Log("get cache from %s, key: %s", groupName, key)

    // 3. 取缓存操作！
    group := GetGroup(groupName)
    if group == nil {
        http.Error(w, "no such group: "+groupName, http.StatusNotFound)
        return
    }

    value, err := group.Get(key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    body, err := proto.Marshal(&pb.Response{Value: value.ByteSlice()})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 4. 返回客户端
    w.Header().Set("Content-Type/octet", "application/octet-stream")
    _, _ = w.Write(body)

}

func (p *HTTPPool) Set(peers ...string) {
    p.mu.Lock()
    defer p.mu.Unlock()

    p.peers = NewMap(defaultReplicas, nil)
    p.peers.Add(peers...)
    p.httpGetters = make(map[string]*httpGetter, len(peers))
    for _, peer := range peers {
        p.httpGetters[peer] = &httpGetter{baseUrl: peer+p.basePath}
    }
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if peer := p.peers.Get(key); peer != "" && peer != p.self {
        p.Log("Pick peer %s", peer)
        return p.httpGetters[peer], true
    }

    return nil, false
}

/*
确保HTTPPool实现PeerPicker接口，如果没有实现，会被IDE识别出来或者在编译的时候报错，
不会在使用过程中才发现没有实现接口
*/
var _ PeerPicker = (*HTTPPool)(nil)