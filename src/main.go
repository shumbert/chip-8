package main

// add package description

import (
	"fmt"
	"os"
)

func main() {
	if i := len(os.Args); i != 2 {
		fmt.Println("Missing argument")
		os.Exit(1)
	}
	program := os.Args[1]

	ioInit()
	machineInitialize()
	machineLoadProgram(program)
	buzz := make(chan struct{})
	draw := make(chan struct{})

	go ioRunBuzzer(buzz)
	go ioRunDisplay(draw)
	go ioRunKeyboard()
	cliRun(buzz, draw)
}
