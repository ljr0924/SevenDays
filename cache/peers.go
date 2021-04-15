package cache

import (
    "fmt"
    "google.golang.org/protobuf/proto"
    "io/ioutil"
    "net/http"
    "net/url"

    pb "SevenDays/cache/cache_pb"
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
    Get(in *pb.Request, out *pb.Response) error
}

type httpGetter struct {
    baseUrl string
}

func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {

    u := fmt.Sprintf(
        "%v/%v/%v",
        h.baseUrl,
        url.QueryEscape(in.GetGroup()),
        url.QueryEscape(in.GetKey()),
    )

    res, err := http.Get(u)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("server returned: %s", res.Status)
    }

    // 读取二进制流
    bytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return fmt.Errorf("reading response body err: %v", err)
    }

    if err = proto.Unmarshal(bytes, out); err != nil {
        return fmt.Errorf("decoding response body: %v", err)
    }

    return nil

}
