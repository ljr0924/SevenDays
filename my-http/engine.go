package my_http

import (
    "net/http"
)

type Engine struct {
    r *router
}

func NewEngine() *Engine {
    return &Engine{
        r: newRouter(),
    }
}

func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
    e.r.addRoute(method, pattern, handler)
}

func (e *Engine) Get(pattern string, handler HandlerFunc) {
    e.addRoute("GET", pattern, handler)
}

func (e *Engine) Post(pattern string, handler HandlerFunc) {
    e.addRoute("POST", pattern, handler)
}

func (e *Engine) Put(pattern string, handler HandlerFunc) {
    e.addRoute("PUT", pattern, handler)
}

func (e *Engine) Delete(pattern string, handler HandlerFunc) {
    e.addRoute("DELETE", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c := newContext(w, r)
    e.r.handle(c)
}

func (e *Engine) Run(addr string) error {
    return http.ListenAndServe(addr, e)
}
