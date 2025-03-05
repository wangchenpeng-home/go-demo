package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from Goroutine!")
}

func main() {
	go sayHello()
	time.Sleep(1 * time.Second)
	fmt.Println("Main function finished!")
}
