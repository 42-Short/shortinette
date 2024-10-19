package main

import (
	"fmt"

	"github.com/42-Short/shortinette/server"
)

func main() {
	r := server.Router()
	if err := r.Run("0.0.0.0:5000"); err != nil {
		fmt.Printf("error running gin server: %v", err)
	}
}
