package main

import "fmt"

func main() {
  name := ""
  fmt.Print("Enter your name: ")
  fmt.Scanf("%s\n", &name)
  fmt.Printf("Hello, %s", name)
}