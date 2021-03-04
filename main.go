package main

import (
    "fmt"
    my_http "my-http"
    "net/http"
)

func main() {
    engine := my_http.NewEngine()
    engine.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "hello %s", r.URL.Query().Get("name"))
    })
    engine.Run(":8080")
}