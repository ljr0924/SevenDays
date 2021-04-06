package cache

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

/*
根据传入的key获取对应的节点
*/
type PeerPicker interface {
    PickPeer(key string) (peer PeerGetter, ok bool)
}

/*
用于从对应的group查找缓存值
*/
type PeerGetter interface {
    Get(group string, key string) ([]byte, error)
}

type httpGetter struct {
    baseUrl string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {

    u := fmt.Sprintf(
        "%v/%v/%v",
        h.baseUrl,
        url.QueryEscape(group),
        url.QueryEscape(group),
    )

    res, err := http.Get(u)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("server returned: %s", res.Status)
    }

    // 读取二进制流
    bytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, fmt.Errorf("reading response body err: %v", err)
    }

    return bytes, nil

}
