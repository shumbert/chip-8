package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

const (
    PROMPT = "chip> "
)

func cliAddBreakpoint() {
    fmt.Println("cliAddBreakpoint")
}

func cliClearBreakpoints() {
    fmt.Println("cliClearBreakpoints")
}

func cliContinueMachine() {
    fmt.Println("cliContinueMachine")
}

func cliDeleteBreakpoint() {
    fmt.Println("cliDeleteBreakpoing")
}

func cliDisassemble(base uint16, count int) {
    for i := 0; i < count; i++ {
        address := base + uint16(i * 2)
        if address < 0x1000 {
            printInstruction(address)
        }
    }
}

func cliExit() {
    fmt.Println()
    os.Exit(0)
}

func cliKillMachine() {
    fmt.Println("cliKillMachine")
}

func cliLoadProgram() {
    fmt.Println("cliLoadProgram")
}

func cliRunMachine() {
    fmt.Println("cliRunMachine")
}

func cliShowBreakpoints() {
    fmt.Println("cliShowBreakpoints")
}

func cliShowHelp() {
    fmt.Println(`Available commands:
e[xit] or q[uit]                quit the interpreter
h[elp]                          show this message

l[oad] <file>                   reset the memory and load a program from file
run                             reset registers and run the machine
s[tep]                          step machine execution
k[ill]                          stop machine execution
c[ontinue]                      resume machine execution

d[isassemble]                   disassemble the next 10 instructions
d[isassemble] <count>           disassemble the next count instructions
d[isassemble] <address> <count> disassemble the next count instructions, starting at address

r[egs]                          show registers
p[ixmap]                        show the display pixmap

b[reak] <address>               set a new breakpoint at address
b[reak]p[oints]                 show breakpoints
del[ete] <breakpoint#>          remove breakpoint number #
cl[ear]                         delete all breakpoints`)
}

func cliShowPixmap() {
    for y := 0; y < SCREENHEIGHT; y++ {
        for x := 0; x < SCREENWIDTH; x++ {
            fmt.Printf("%d", m.pixmap[x][y])
        }
        fmt.Printf("\n")
    }
}

func cliShowRegs() {
    fmt.Printf("[V%X] 0x%02x    [DT]=0x%02x    [SP%X] 0x%03x\n", 0x0, m.regs.v[0x0], m.regs.dt, 0x0, m.stack[0x0])
    fmt.Printf("[V%X] 0x%02x    [ST]=0x%02x    [SP%X] 0x%03x\n", 0x1, m.regs.v[0x1], m.regs.st, 0x1, m.stack[0x1])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0x2, m.regs.v[0x2], 0x2, m.stack[0x2])
    fmt.Printf("[V%X] 0x%02x    [I]=0x%03x    [SP%X] 0x%03x\n", 0x3, m.regs.v[0x3], m.regs.i, 0x3, m.stack[0x3])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0x4, m.regs.v[0x4], 0x4, m.stack[0x4])
    fmt.Printf("[V%X] 0x%02x    [PC]=0x%03x   [SP%X] 0x%03x\n", 0x5, m.regs.v[0x5], m.regs.pc, 0x5, m.stack[0x5])
    fmt.Printf("[V%X] 0x%02x    [SP]=0x%02x    [SP%X] 0x%03x\n", 0x6, m.regs.v[0x6], m.regs.sp, 0x6, m.stack[0x6])
    fmt.Printf("[V%X] 0x%02x    [CY]=0x%01x     [SP%X] 0x%03x\n", 0x7, m.regs.v[0x7], m.cycles, 0x7, m.stack[0x7])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0x8, m.regs.v[0x8], 0x8, m.stack[0x8])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0x9, m.regs.v[0x9], 0x9, m.stack[0x9])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0xa, m.regs.v[0xa], 0xa, m.stack[0xa])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0xb, m.regs.v[0xb], 0xb, m.stack[0xb])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0xc, m.regs.v[0xc], 0xc, m.stack[0xc])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0xd, m.regs.v[0xd], 0xd, m.stack[0xd])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0xe, m.regs.v[0xe], 0xe, m.stack[0xe])
    fmt.Printf("[V%X] 0x%02x                 [SP%X] 0x%03x\n", 0xf, m.regs.v[0xf], 0xf, m.stack[0xf])
}

func printInstruction(address uint16) {
	assembled := getInstruction(address)
    disassembled := disassembleInstruction(assembled)

	fmt.Printf("0x%03x: 0x%04x ", address, assembled)
    switch {
    case disassembled.op == sys:
		fmt.Printf("SYS")

    case disassembled.op == cls:
		fmt.Printf("SYS")

    case disassembled.op == ret:
		fmt.Printf("RET")

    case disassembled.op == jmp:
		fmt.Printf("JMP 0x%03x", disassembled.nnn)

    case disassembled.op == call:
		fmt.Printf("CALL 0x%03x", disassembled.nnn)

    case disassembled.op == seb:
		fmt.Printf("SE V%X, 0x%02x", disassembled.x, disassembled.kk)

    case disassembled.op == sneb:
		fmt.Printf("SNE V%X, 0x%02x", disassembled.x, disassembled.kk)

    case disassembled.op == ser:
		fmt.Printf("SE V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == ldb:
		fmt.Printf("LD V%X, 0x%02x", disassembled.x, disassembled.kk)

    case disassembled.op == addb:
		fmt.Printf("ADD V%X, 0x%02x", disassembled.x, disassembled.kk)

    case disassembled.op == ldr:
		fmt.Printf("LD V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == or:
		fmt.Printf("OR V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == and:
		fmt.Printf("AND V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == xor:
		fmt.Printf("XOR V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == addr:
		fmt.Printf("ADD V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == sub:
		fmt.Printf("SUB V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == shr:
		fmt.Printf("SHR V%X {, V%X}", disassembled.x, disassembled.y)

    case disassembled.op == subn:
		fmt.Printf("SUBN V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == shl:
		fmt.Printf("SHL V%X {, V%X}", disassembled.x, disassembled.y)

    case disassembled.op == sner:
		fmt.Printf("SNE V%X, V%X", disassembled.x, disassembled.y)

    case disassembled.op == ldi:
		fmt.Printf("LD I, 0x%03x", disassembled.nnn)

    case disassembled.op == jpv:
		fmt.Printf("JP V0, 0x%03x", disassembled.nnn)

    case disassembled.op == rnd:
		fmt.Printf("RND V%X, 0x%02x", disassembled.x, disassembled.kk)

    case disassembled.op == drw:
		fmt.Printf("DRW V%X, V%X, 0x%x", disassembled.x, disassembled.y, disassembled.n)

    case disassembled.op == skp:
		fmt.Printf("SKP V%X", disassembled.x)

    case disassembled.op == sknp:
		fmt.Printf("SKNP V%X", disassembled.x)

    case disassembled.op == gett:
		fmt.Printf("LD V%X, DT", disassembled.x)

    case disassembled.op == ldk:
		fmt.Printf("LD V%X, K", disassembled.x)

    case disassembled.op == sett:
		fmt.Printf("LD DT, V%X", disassembled.x)

    case disassembled.op == lds:
		fmt.Printf("LD ST, V%X", disassembled.x)

    case disassembled.op == addi:
		fmt.Printf("ADD I, V%X", disassembled.x)

    case disassembled.op == ldf:
		fmt.Printf("LD F, V%X", disassembled.x)

    case disassembled.op == ldbcd:
		fmt.Printf("LD B, V%X", disassembled.x)

    case disassembled.op == save:
		fmt.Printf("LD [I], V%X", disassembled.x)

    case disassembled.op == restore:
		fmt.Printf("LD V%X, [I]", disassembled.x)
    }
	fmt.Printf("\n")
}

func runCLI() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Type \"h\" or \"help\" for commands usage")

    for {
        fmt.Printf(PROMPT)
        input, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                cliExit()
            }
            fmt.Fprintln(os.Stderr, err)
        }
        input = strings.TrimSuffix(input, "\n")

        // FIXME: command is not properly parsed if there
        // are spaces
        // for instance s;d 1;r
        // but s ; d 1 ; r is not
        for _, command := range strings.Split(input, ";") {
            args := strings.Split(command, " ")

            switch args[0] {
            case "b", "break":
                cliAddBreakpoint()

            case "bp", "breakpoints":
                cliShowBreakpoints()

            case "cl", "clear":
                cliClearBreakpoints()

            case "c", "continue":
                cliContinueMachine()

            case "delete":
                cliDeleteBreakpoint()

            case "d", "disassemble":
                base := m.regs.pc
                count := 10
                if len(args) > 2 {
                    count, _ = strconv.Atoi(args[2])
                    if strings.HasPrefix(args[1], "0x") || strings.HasPrefix(args[2], "0X") {
                        i, _ := strconv.ParseInt(args[1][2:], 16, 16)
                        base = uint16(i)
                    } else {
                        i, _ := strconv.ParseInt(args[1], 10, 16)
                        base = uint16(i)
                    }
                } else if len(args) > 1 {
                    count, _ = strconv.Atoi(args[1])
                }
                cliDisassemble(base, count)

            case "e", "exit":
                cliExit()

            case "h", "help":
                cliShowHelp()

            case "k", "kill":
                cliKillMachine()

            case "l", "load":
                cliLoadProgram()

            case "p", "pixmap":
                cliShowPixmap()

            case "q", "quit":
                cliExit()

            case "r", "regs":
                cliShowRegs()

            case "run":
                //cliRunMachine()
                runMachine()

            case "s", "step":
                stepMachine()

            default:
                fmt.Printf("%s: unrecognized command\n", args[0])
            }
        }
    }
}
