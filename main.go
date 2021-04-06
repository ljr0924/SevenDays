package main

import (
    "SevenDays/cache"
    "flag"
    "fmt"
    "log"
    "net/http"
)

var db = map[string]string {
    "Tom": "630",
    "Jack": "589",
    "Sam": "567",
}

func createGroup() *cache.Group {
    return cache.NewGroup("scores", 2<<10, cache.GetterFunc(func(key string) ([]byte, error) {
        log.Println("[Slow DB] search key: ", key)
        if v, ok := db[key];ok {
            return []byte(v), nil
        }
        return nil, fmt.Errorf("%s not exists", key)
    }))
}

func startCacheServer(addr string, addrList []string, group *cache.Group) {
    peers := cache.NewHTTPPool(addr)
    peers.Set(addrList...)
    group.RegisterPeers(peers)
    log.Println("cache is running at ", addr)
    log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startApiServer(apiAddr string, group *cache.Group) {
    http.Handle("/api", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        key := request.URL.Query().Get("key")
        view, err := group.Get(key)
        if err != nil {
            http.Error(writer, err.Error(), http.StatusInternalServerError)
            return
        }
        writer.Header().Set("Content-Type", "application/octet-stream")
        _, _ = writer.Write(view.ByteSlice())
    }))
    log.Println("frontend server is running at", apiAddr)
    log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
    var port int
    var api bool
    flag.IntVar(&port, "port", 8001, "cache server port")
    flag.BoolVar(&api, "api", false, "Start a api server?")
    flag.Parse()

    apiAddr := "http://localhost:9999"
    addrMap := map[int]string{
        8001: "http://localhost:8001",
        8002: "http://localhost:8002",
        8003: "http://localhost:8003",
    }

    var addrList []string
    for _, v := range addrMap {
        addrList = append(addrList, v)
    }

    group := createGroup()
    if api {
        go startApiServer(apiAddr, group)
    }
    startCacheServer(addrMap[port], []string(addrList), group)
}