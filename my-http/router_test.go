package my_http

import (
    "fmt"
    "reflect"
    "testing"
)

func newTestRouter() *router {
    r := newRouter()
    r.addRoute("GET", "/", nil)
    r.addRoute("GET", "/hello/:name", nil)
    r.addRoute("GET", "/hello/:name/age", nil)
    r.addRoute("GET", "/hello/b/c", nil)
    // r.addRoute("GET", "/hi/:name", nil)
    // r.addRoute("GET", "/assets/*filepath", nil)
    return r
}


func TestParsePattern(t *testing.T) {
    ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
    ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
    ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
    if !ok {
        t.Fatal("test parsePattern failed")
    }
}

func TestGetRoute(t *testing.T) {
    r := newTestRouter()
    n, ps := r.getRoute("GET", "/hello/banana/age")

    if n == nil {
        t.Fatal("nil shouldn't be returned")
    }

    if n.Pattern != "/hello/:name/age" {
        t.Fatal("should match /hello/:name")
    }

    if ps["name"] != "banana" {
        t.Fatal("name should be equal to 'banana'")
    }

    fmt.Printf("matched path: %s, params['name']: %s\n", n.Pattern, ps["name"])

}