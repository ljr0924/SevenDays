package cache

import (
    "fmt"
    "log"
    "net/http"
    "testing"
)


func TestHttpPool(t *testing.T) {
    NewGroup("scores", 2<<10, GetterFunc(
        func(key string) ([]byte, error) {
            log.Println("[SlowDB] search key", key)
            if v, ok := db[key]; ok {
                return []byte(v), nil
            }
            return nil, fmt.Errorf("%s not exist", key)
        }))

    addr := "localhost:8080"
    peers := NewHTTPPool(addr)
    log.Println("running at", addr)
    log.Fatal(http.ListenAndServe(addr, peers))
}
