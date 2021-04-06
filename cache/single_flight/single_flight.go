package single_flight

import "sync"

/* 代表正在进行中，或已经结束的请求，使用sync.WaitGroup避免重入 */
type call struct {
    wg  sync.WaitGroup
    val interface{}
    err error
}

/* 管理不同key的请求 */
type Group struct {
    mu sync.Mutex
    m  map[string]*call
}

/*
针对相同的key，无论调用多少次Do，函数fn都只会调用一次，等待fn调用结束，才会返回值或错误
*/
func (g *Group) Do(key string, fn func()(interface{}, error)) (interface{}, error) {
    g.mu.Lock()
    if g.m == nil {
        g.m = make(map[string]*call)
    }

    if c, ok := g.m[key]; ok {
        g.mu.Unlock()
        c.wg.Wait()    // 如果请求正在进行中，则等待
        return c.val, c.err // 请求结束，返回结果
    }

    c := new(call)
    c.wg.Add(1)    // 向WaitGroup加锁
    g.m[key] = c         // 添加到g.m表明正在进行请求
    g.mu.Unlock()

    c.val, c.err = fn()  // 调用fn，发起请求
    c.wg.Done()          // 请求结束

    g.mu.Lock()
    delete(g.m, key)     // 更新g.m
    g.mu.Unlock()

    return c.val, c.err

}