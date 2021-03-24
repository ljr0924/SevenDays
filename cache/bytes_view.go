package cache

type ByteView struct {
    b []byte  // 真实的缓存值。选择byte类型是为了能够支持缓存任意数据类型，例如字符串，图片等
}

// 返回缓存值的大小
func (bv ByteView) Len() int {
    return len(bv.b)
}

// b是只读的，使用ByteSlice方法返回一个拷贝，防止缓存值被外部程序修改
func (bv ByteView) ByteSlice() []byte {
    return cloneBytes(bv.b)
}

func (bv ByteView) String() string {
    return string(bv.b)
}

func cloneBytes(b []byte) []byte {
    c := make([]byte, len(b))
    copy(c, b)
    return c
}

