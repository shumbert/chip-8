package main

// add package description

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    if i := len(os.Args); i != 2 {
	    fmt.Println("Missing argument")
        os.Exit(1)
    }
    program := os.Args[1]

    go runDisplay()
	resetMachine()
    loadProgram(program)
    printMachineState()
    printDisplay()

    buf := bufio.NewReader(os.Stdin)
    for {
        buf.ReadBytes('\n')
        fmt.Println("----------------")
        printInstruction(state.pc)
        stepMachine()
        printMachineState()
        printDisplay()
    }
}
