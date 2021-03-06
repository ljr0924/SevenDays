package cache

import (
    "SevenDays/cache/container"
    "reflect"
    "testing"
)

type String string

func (d String) Len() int {
    return len(d)
}

func TestGet(t *testing.T) {
    cache := container.NewLru(int64(0), nil)
    cache.Set("key1", String("1234"))
    if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "1234" {
        t.Fatalf("cache hit key1=1234 failed")
    }
    if _, ok := cache.Get("key2"); ok {
        t.Fatalf("cache miss key2 failed")
    }
}

func TestRemoveOldest(t *testing.T) {
    k1, k2, k3 := "key1", "key2", "k3"
    v1, v2, v3 := "value1", "value2", "v3"
    cap := len(k1 + k2 + v1 + v2)
    cache := container.NewLru(int64(cap), nil)
    cache.Set(k1, String(v1))
    cache.Set(k2, String(v2))
    cache.Set(k3, String(v3))

    if _, ok := cache.Get("key1"); ok || cache.Len() != 2 {
        t.Fatalf("Removeoldest key1 failed")
    }
}

func TestOnEvictedFunc(t *testing.T) {

    keys := make([]string, 0)

    callback := func(key string, value container.Value) {
        keys = append(keys, key)
    }

    cache := container.NewLru(8, callback)
    cache.Set("k1", String("v1"))
    cache.Set("k2", String("v2"))
    cache.Set("k3", String("v3"))
    cache.Set("k4", String("v4"))

    except := []string{"k1", "k2"}
    if !reflect.DeepEqual(except, keys) {
        t.Fatalf("keys not equal except, except: %v  keys: %v", except, keys)
    }

    t.Logf("keys equal except, except: %v  keys: %v", except, keys)
}

