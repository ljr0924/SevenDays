package my_http

import (
    "fmt"
    "testing"
)

func TestEngineServeHTTP(t *testing.T) {

    engine := NewEngine()
    engine.GET("/json", func(c *Context) {
        c.JSON(200, H{
            "name": "msh",
            "age": 19,
            "feature": "sb",
        })
    })
    engine.GET("/data", func(c *Context) {
        c.Data(200, []byte("这是data"))
    })
    engine.GET("/string", func(c *Context) {
        c.String(200, "这是string")
    })
    engine.GET("/html", func(c *Context) {
        c.HTML(200, "<h1>这是html</h1>")
    })
    engine.GET("/hello/:name", func(c *Context) {
        c.HTML(200, fmt.Sprintf("<h1>hello %s</h1>", c.Param("name")))
    })
    engine.Run(":8080")

}

