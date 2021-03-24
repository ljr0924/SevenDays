package container

import "container/list"

type Lru struct {
    maxBytes int64  // 允许使用最大内存
    nBytes   int64  // 已使用内存
    ll       *list.List
    // 字典的定义是 map[string]*list.Element，键是字符串，值是双向链表中对应节点的指针。
    cache    map[string]*list.Element

    // 某条记录被移除时的回调函数，可以为 nil
    onEvicted func(key string, value Value)
}

func NewLru(maxBytes int64, onEvicted func(key string, value Value)) *Lru {
    return &Lru{
        maxBytes: maxBytes,
        nBytes: 0,
        ll: list.New(),
        cache: make(map[string]*list.Element),
        onEvicted: onEvicted,
    }
}

type entry struct {
    key   string
    value Value
}

func (e *entry) getBytes() int64 {
    return int64(len(e.key)) + int64(e.value.Len())
}

func (lru *Lru) Get(key string) (Value, bool) {
    if ele, ok := lru.cache[key]; ok {
        lru.ll.MoveToFront(ele)
        kv := ele.Value.(*entry)
        return kv.value, true
    }

    return nil, false
}

func (lru *Lru) RemoveOldest() {
    // 1. 找到队首元素
    ele := lru.ll.Back()
    if ele != nil {
        // 2. 从队列中移除元素
        lru.ll.Remove(ele)
        kv := ele.Value.(*entry)
        // 3. 从缓存里移除数据
        delete(lru.cache, kv.key)
        // 4. 更新cache已使用内存数
        lru.nBytes -= kv.getBytes()
        // 5. 如果有回调函数，调用回调函数
        if lru.onEvicted != nil {
            lru.onEvicted(kv.key, kv.value)
        }
    }
}

func (lru *Lru) Add(key string, value Value) {
    if ele, ok := lru.cache[key]; ok {
        lru.ll.MoveToFront(ele)
        kv := ele.Value.(*entry)
        lru.nBytes += int64(value.Len()) - int64(kv.value.Len())
        kv.value = value
    } else {
        ele := lru.ll.PushFront(&entry{key, value})
        lru.cache[key] = ele
        lru.nBytes += int64(len(key)) + int64(value.Len())
    }
    // 删除旧缓存
    for lru.maxBytes != 0 && lru.maxBytes < lru.nBytes {
        lru.RemoveOldest()
    }
}

func (lru *Lru) Len() int {
    return lru.ll.Len()
}
