package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "log"
    "errors"
    "io"
    "github.com/google/uuid"
)

type task struct {
    Completed bool `json:"completed"`
    ID string `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
}

var tasks map[string]task = make(map[string]task)

func taskHandler(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var re, result task
    err := decoder.Decode(&re)

    if err != nil && !errors.Is(err, io.EOF) {
        panic(err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if req.Method == "GET" {
        encoder := json.NewEncoder(w)
        encoder.Encode(tasks)
    }
    if req.Method == "POST" || req.Method == "PUT" {
	i := re.ID
	if i == "" {
		i = fmt.Sprintf("%s", uuid.New()) 
		re.ID = i
	}
        tasks[i] = re
        result = tasks[i]
        encoder := json.NewEncoder(w)
        encoder.Encode(&result)
    }
    if req.Method == "DELETE" {
	delete(tasks, re.ID)
    }
}

func main() {
    http.HandleFunc("/tasks", taskHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
