# Lab GO01 - A Simple Go App

## Objective
Create and run a simple application using golang

## Outcomes
* Gain exposure to core language constructs in go
* Build and run go apps

## High-Level Steps
* Create a simple Golang app
* Execute an app with go run
* Compile an app with go build

## Detailed Steps

### Ensure Go is installed
We will begin by ensuring that Go is installed. Run the following command:
```bash
go version
```
Go 1.25 should already be installed in the lab environment

### Hello World
Change directory into the lab01 folder:
```bash
cd go-labs/lab01
```
Open main.go in the editor
In the editor, add the following to main.go:
```go
package main

import "fmt"

func main() {
  fmt.Println("Hello World!")
}
```
In your terminal, execute the file using _go run_:
```bash
go run main.go
```
You should see the message 'Hello World!' displayed in the terminal.

### Handle User Input
Edit the main.go file again, and amend the contents as follows:
```go
package main

import "fmt"

func main() {
  name := ""
  fmt.Print("Enter your name: ")
  fmt.Scanf("%s\n", &name)
  fmt.Printf("Hello, %s", name)
}
```
The extra lines of code added here initialise an empty string variable, before displaying a prompt to the user and reading input until a newline is reached, storing the input into the name variable. A message which includes the name is then printed. 

Execute the file again:
```bash
go run main.go
```

### Compile the Application
Instead of using the _go run_ command to interpret the code each time we run the app, we can instead compile it, for faster execution. We will use _go build_ for this. First, we need to create a minimal _go.mod_ file, as go build will not work without one:
```bash
echo "module greeter" > go.mod
```
A module name is the minimum configuration required for go build to be able to build the app. To compile the code, run the following:
```bash
go build -o greeter
```
This will produce an executable called _greeter_. Execute the file from the terminal:
```bash
./greeter
```
