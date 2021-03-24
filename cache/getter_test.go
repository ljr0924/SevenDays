package cache

import (
    "fmt"
    "reflect"
    "testing"
)

func TestGetterFunc(t *testing.T) {

    var f Getter = GetterFunc(func(key string) ([]byte, error) {
        return []byte(key), nil
    })

    buf, _ := f.Get("123")
    fmt.Println(buf)

    expect := []byte("123")
    if v, _ := f.Get("123"); !reflect.DeepEqual(v, expect) {
        t.Errorf("callback failed")
    }

}
