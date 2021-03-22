package my_http

import (
    "log"
    "net/http"
    "strings"
)

type Engine struct {
    *RouterGroup
    groups []*RouterGroup
    r      *router
}

type RouterGroup struct {
    prefix     string        // 前缀
    middleware []HandlerFunc // 中间件，支持扩展功能
    parent     *RouterGroup  // 父节点
    engine     *Engine       //
}

func NewEngine() *Engine {
    engine := &Engine{r: newRouter()}
    engine.RouterGroup = &RouterGroup{engine: engine}
    engine.groups = make([]*RouterGroup, 0)
    return engine
}

func (rg *RouterGroup) Use(middleware ...HandlerFunc) {
    rg.middleware = append(rg.middleware, middleware...)
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
    engine := rg.engine
    newGroup := &RouterGroup{
        prefix: rg.prefix + prefix,
        parent: rg,
        engine: engine,
    }
    engine.groups = append(engine.groups, newGroup)
    return newGroup
}

func (rg *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
    pattern := rg.prefix + comp
    rg.engine.r.addRoute(method, pattern, handler)
}

func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
    rg.addRoute("GET", pattern, handler)
}

func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
    rg.addRoute("POST", pattern, handler)
}

func (rg *RouterGroup) PUT(pattern string, handler HandlerFunc) {
    rg.addRoute("PUT", pattern, handler)
}

func (rg *Engine) DELETE(pattern string, handler HandlerFunc) {
    rg.addRoute("DELETE", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var middleware []HandlerFunc
    middleware = append(middleware, e.middleware...)
    for _, group := range e.groups {
        if strings.HasPrefix(r.URL.Path, group.prefix) {
            middleware = append(middleware, group.middleware...)
        }
    }
    c := newContext(w, r)
    c.handlers = middleware
    e.r.handle(c)
}

func (e *Engine) Run(addr string) error {
    var host string
    if strings.HasPrefix(addr, ":") {
        host = "http://localhost" + addr
    }
    log.Printf("server run at ==> %s", host)
    return http.ListenAndServe(addr, e)
}
