package cache

import (
    "math/rand"
    "testing"
    "time"
)

var db = map[string]string{
    "123": "123",
    "456": "456",
    "789": "789",
}

func GetFromMap(key string) string {
    time.Sleep(3 * time.Second)
    return db[key]
}

func TestGroup(t *testing.T) {

    keys := []string{"123", "456", "789", "111", "222"}

    group := NewGroup("cache1", 12, GetterFunc(func(key string) ([]byte, error) {
        t.Logf("get from db by key : %s", key)
        v := GetFromMap(key)
        if v != "" {
            return []byte(key), nil
        }
        return []byte{}, nil
    }))

    for i := 0; i < 20; i++ {
        n := rand.Intn(4)
        k := keys[n]

        v, _ := group.Get(k)
        t.Log(v.String())
        t.Log("--------------------")

    }


}
