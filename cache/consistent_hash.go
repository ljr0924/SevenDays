package cache

import (
    "hash/crc32"
    "sort"
    "strconv"
)

type Hash func(data []byte) uint32

// 包含所有哈希键
type Map struct {
    hash     Hash
    replicas int // 虚拟节点数目
    keys     []int // hash环 sorted
    hashMap  map[int]string // 虚拟节点与真实节点的映射表  K:虚拟节点 V:真实节点
}

/*
采取依赖注入的方式，允许用户换成自定义的Hash函数
默认为crc32.ChecksumIEEE
*/
func NewMap(replicas int, fn Hash) *Map {
    m := &Map{
        hash:     fn,
        replicas: replicas,
        hashMap:  make(map[int]string),
    }
    if m.hash == nil {
        m.hash = crc32.ChecksumIEEE
    }

    return m
}

func (m *Map) Add(keys ...string) {
    for _, key := range keys {
        for i := 0; i < m.replicas; i++ {
            // 计算key的
            hash := int(m.hash([]byte(strconv.Itoa(i)+key)))
            // 加到hash
            m.keys = append(m.keys, hash)
            // 添加虚拟节点与真实节点映射
            m.hashMap[hash] = key
        }
    }
    // 排序
    sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
    if len(m.keys) == 0 {
        return ""
    }

    hash := int(m.hash([]byte(key)))

    // 二分查找
    idx := sort.Search(len(m.keys), func(i int) bool {
        return m.keys[i] >= hash
    })

    return m.hashMap[m.keys[idx%len(m.keys)]]
}



