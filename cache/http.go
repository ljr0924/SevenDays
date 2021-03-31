package cache

import (
    "log"
    "net/http"
    "strings"
)

const defaultBasePath = "/_banana"

type HTTPPool struct {
    self     string  // 用来记录自己的地址，包括主机名/IP 和端口。
    basePath string
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

    w.Header().Set("Content-Type/octet", "application/octet-stream")
    _, _ = w.Write(value.ByteSlice())

}
