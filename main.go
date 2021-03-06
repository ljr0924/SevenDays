package main

import (
    my_http "my-http"
)

func main() {
    engine := my_http.NewEngine()
    engine.GET("/hello", func(c *my_http.Context) {
        c.String(200, "hello")
    })
    engine.Run(":8080")
}