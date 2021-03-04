package my_http

import "fmt"

// 自定义处理函数类型
type HandlerFunc func(*Context)

type router struct {
    handlers map[string]HandlerFunc
}

func newRouter() *router {
    return &router{
        handlers: map[string]HandlerFunc{},
    }
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
    key := r.getRouteKey(method, pattern)
    r.handlers[key] = handler
}

func (r *router) getRouteKey(method, pattern string) string {
    return method + "-" + pattern
}

func (r *router) handle(c *Context) {
    key := r.getRouteKey(c.Method, c.Path)
    if handler, ok := r.handlers[key]; ok {
        handler(c)
    } else {
        fmt.Fprintf(c.w, "404 Not Found: %q", c.r.URL)
    }
}
