# Lab GO02 - Create a Web App in Go

## Objective
Create and test a reasonably simple web application using Golang

## Outcomes
* Use the go std library's http module to create a web app
* Write unit tests using the testing module

## High-Level Steps
* Write web app code
* Manually test the web app using curl
* Implement unit tests for the app in golang

## Detailed Steps
### Implement the Web App
In the terminal, change directory into the lab02 folder. Open the main.go file in the editor:
```go
package main

import (
    "fmt"
)

func main() {
    fmt.Println("TODO: Implement web app");
}
```
We will need to make some changes to this file, starting with the imports. Since we are creating a web app we will need to be able to handle HTTP requests - we can use the standard library net/http module for this. We will also need a few other utility modules. The imports should look like:
```go
import (
    "encoding/json"
    "net/http"
    "log"
    "errors"
    "io"
)
```
Note that we are no longer importing the fmt library - we will not need it for the web app. Next, add the following between the import block and the main function:
```go
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
        tasks[re.ID] = re
        result = tasks[re.ID]
        encoder := json.NewEncoder(w)
        encoder.Encode(&result)
    }
    if req.Method == "DELETE" {
        delete(tasks, re.ID)
    }
}
```
This is the core logic of our web app. Listening for http requests, it will parse JSON from the request body according to the schema represented by the struct type definition. Now we need to edit the main function to actually serve the application:
```go
func main() {
    http.HandleFunc("/tasks", taskHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
This registers the function containing our application logic as a handler for requests to the _/calc_ endpoint, before starting the server with logging enabled.

### Build and Run the Application
Compile the code using go build:
```bash
go build -o bin/webapp
```
then make the binary executable, and run it:
```bash
chmod +x bin/webapp
./bin/webapp
```
Once the app is running, in a new terminal session use curl to test the functionality of the app:
```bash
curl -XPOST -d@request.json http://localhost:8080/tasks
curl -XGET http://localhost:8080/tasks
# edit request.json, setting 'Completed' to true
curl -XPUT -d@request.json http://localhost:8080/tasks
curl -XGET http://localhost:8080/tasks
curl -XDELETE -d@request.json http://localhost:8080/tasks
curl -XGET http://localhost:8080/tasks
```
In the original terminal, stop the app with ctrl+c

### Add a Basic Unit Test
We will now add a basic unit test for our app. For this, we will need to use an external library to _mock_ the response from the app. Run the following:
```bash
go get -u github.com/h2non/gock
```
Once the gock package has been installed, create a new file in the same folder called _main\_test.go_, with the following contents:
```go
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
```
Execute the test with _go test_:
```bash
go test
```
You should observe a result similar to:
```text
PASS
ok      app     0.004s
```
