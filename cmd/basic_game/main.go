package main

import (
	"log"
)

func main() {
	svr := NewServer()
	log.Fatal(svr.Start(8080))
}
