package container

type Value interface {
    Len() int  // 用于返回值所占用的内存大小
}

type Container interface {
    Get(key string) (Value, bool)
    Set(key string, value Value)
    Len() int
    RemoveOldest()
}
