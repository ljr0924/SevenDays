package my_http

import (
    "testing"
)

func TestEngineServeHTTP(t *testing.T) {

    engine := NewEngine()
    engine.Get("/json", func(c *Context) {
        c.JSON(200, H{
            "name": "msh",
            "age": 19,
            "feature": "sb",
        })
    })
    engine.Get("/data", func(c *Context) {
        c.Data(200, []byte("这是data"))
    })
    engine.Get("/string", func(c *Context) {
        c.String(200, "这是string")
    })
    engine.Get("/html", func(c *Context) {
        c.HTML(200, "<h1>这是html</h1>")
    })
    engine.Run(":8080")

}

