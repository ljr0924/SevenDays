package my_http

import (
    "log"
    "testing"
    "time"
)

func v2Middleware() HandlerFunc {
    return func(c *Context) {
        log.Println("this is v2 start point!!!!")
        c.Next()
        log.Println("this is v2 end point!!!!")
    }
}

func TestMiddleware(t *testing.T) {

    r := NewEngine()
    r.Use(Logger())

    v1 := r.Group("/v1")
    v1.GET("/hello", func(context *Context) {
        time.Sleep(time.Second * 3)
        context.JSON(200, map[string]string{
            "name": "梁嘉荣",
            "age":  "18",
        })
    })

    v2 := r.Group("/v2")
    v2.Use(v2Middleware())
    v2.GET("/hello", func(context *Context) {
        time.Sleep(time.Second * 3)
        context.JSON(200, map[string]string{
            "name": "梁嘉荣",
            "age":  "18",
        })
    })

    r.Run(":8080")
}
