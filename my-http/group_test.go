package my_http

import "testing"

func TestGroup(t *testing.T) {

    e := NewEngine()

    handler := func(context *Context) {
        context.HTML(200, "<h1>" + context.r.URL.Path + "</h1>")
    }

    e.GET("/index", handler)

    sub := e.Group("/index/3")
    sub.GET("/4", handler)

    e.Run(":8080")

}
