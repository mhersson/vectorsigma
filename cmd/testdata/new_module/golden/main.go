package main

import (
	"fmt"

	"new_module/internal/fsm"
)

func main() {
	SM := fsm.New()
	err := SM.Run()
	if err != nil {
		fmt.Printf("State machine run ended with error: %v\n", err)
	}
}
