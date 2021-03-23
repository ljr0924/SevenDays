package my_http

import "testing"

func TestStatic(t *testing.T) {

    e := NewEngine()
    e.Static("/assets", "/Users/banana/git_project/SevenDays/my-http/static")

    e.Run(":8080")
}
