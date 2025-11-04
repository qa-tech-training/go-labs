package main

import (
    "net/http"
    "testing"
    "github.com/h2non/gock"
)

func TestTaskHandler(t *testing.T) {
    defer gock.Off()

    p := task{ID: "abcd", Title: "Foo", Description: "Bar", Completed: false}
    gock.New("http://localhost:8080").
        Get("/tasks").
        Reply(200).
        JSON(p)

    res, err := http.Get("http://localhost:8080/tasks")
    if err != nil || res.StatusCode != 200 {
        t.Errorf("Get Request Failed")
    }
}
