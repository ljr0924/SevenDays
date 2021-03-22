package my_http

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type H map[string]interface{}

type Context struct {

    // 原始对象
    w http.ResponseWriter
    r *http.Request

    // 请求相关
    Method string
    Path   string
    Params map[string]string

    // 响应相关
    StatusCode int

    // 中间件
    handlers []HandlerFunc
    index    int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
    return &Context{
        w:      w,
        r:      r,
        Method: r.Method,
        Path:   r.URL.Path,
        index:  -1,
    }
}

func (c *Context) Next() {
    c.index++
    s := len(c.handlers)
    for ; c.index < s; c.index++ {
        c.handlers[c.index](c)
    }
}

func (c *Context) Param(key string) string {
    v, _ := c.Params[key]
    return v
}

func (c *Context) PostForm(key string) string {
    return c.r.PostFormValue(key)
}

func (c *Context) Query(key string) string {
    return c.r.URL.Query().Get(key)
}

func (c *Context) SetHeader(key, value string) {
    c.r.Header.Set(key, value)
}

func (c *Context) Status(code int) {
    c.w.WriteHeader(code)
}

func (c *Context) String(code int, format string, values ...interface{}) {
    c.SetHeader("Content-Type", "text/plain")
    _, _ = fmt.Fprintf(c.w, format, values...)
    c.Status(code)
}

func (c *Context) JSON(code int, obj interface{}) {
    c.SetHeader("Content-Type", "application/json")
    c.Status(code)
    encoder := json.NewEncoder(c.w)
    if err := encoder.Encode(obj); err != nil {
        http.Error(c.w, "error system", 500)
    }
}

func (c *Context) Data(code int, data []byte) {
    c.Status(code)
    c.w.Write(data)
}

func (c *Context) HTML(code int, html string) {
    c.SetHeader("Content-Type", "text/html")
    c.Status(code)
    c.w.Write([]byte(html))
}
