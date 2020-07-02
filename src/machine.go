package main

// add package description

import (
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "time"
)

// define and declare structs for:
// memory
// state (registers)
// stack

// type byte is synonymous to uint8
// bytes package

const (
    MEMFONTS = 0x000
    MEMPROGRAMSTART = 0x200
    MEMEND = 0x1000
    SCREENWIDTH = 64
    SCREENHEIGHT = 32
    SLEEPTIME = 16666667 * time.Nanosecond
)

type opcode int
const (
    sys opcode = iota // 0nnn - SYS addr
	cls               // 00E0 - CLS
	ret               // 00EE - RET
	jmp               // 1nnn - JP addr
	call              // 2nnn - CALL addr
	seb               // 3xkk - SE Vx, byte
	sneb              // 4xkk - SNE Vx, byte
	ser               // 5xy0 - SE Vx, Vy
	ldb               // 6xkk - LD Vx, byte
	addb              // 7xkk - ADD Vx, byte
	ldr               // 8xy0 - LD Vx, Vy
	or                // 8xy1 - OR Vx, Vy
	and               // 8xy2 - AND Vx, Vy
	xor               // 8xy3 - XOR Vx, Vy
	addr              // 8xy4 - ADD Vx, Vy
	sub               // 8xy5 - SUB Vx, Vy
	shr               // 8xy6 - SHR Vx {, Vy}
	subn              // 8xy7 - SUBN Vx, Vy
	shl               // 8xyE - SHL Vx {, Vy}
	sner              // 9xy0 - SNE Vx, Vy
	ldi               // Annn - LD I, addr
	jpv               // Bnnn - JP V0, addr
	rnd               // Cxkk - RND Vx, byte
	drw               // Dxyn - DRW Vx, Vy, nibble
	skp               // Ex9E - SKP Vx
	sknp              // ExA1 - SKNP Vx
	gett              // Fx07 - LD Vx, DT
	ldk               // Fx0A - LD Vx, K
	sett              // Fx15 - LD DT, Vx
	lds               // Fx18 - LD ST, Vx
	addi              // Fx1E - ADD I, Vx
	ldf               // Fx29 - LD F, Vx
	ldbcd             // Fx33 - LD B, Vx
	save              // Fx55 - LD [I], Vx
	restore           // Fx65 - LD Vx, [I]
)

type instruction struct {
	op opcode  // Instruction opcode
	nnn uint16 // (or addr) A 12-bit value, the lowest 12 bits of the instruction
	x byte     // A 4-bit value, the lower 4 bits of the high byte of the instruction
	y byte     // A 4-bit value, the upper 4 bits of the low byte of the instruction
	kk byte    // (or byte) An 8-bit value, the lowest 8 bits of the instruction
	n byte     // (or nibble) A 4-bit value, the lowest 4 bits of the instruction
}

// 
// Definition of the machine state
//

// First the machine registers
type registers struct {
	v  [16]byte // Data registers V0 to VF
	i  uint16   // Address Register
	dt byte     // Delay Timer
	st byte     // Sound Timer
	pc uint16   // Program Counter
	sp byte     // Stack Pointer
}

// Now let's bundle the register with the machine memory,
// the stack and the display pixmap.
//
// Hardware running CHIP-8 were typically clocked at 540Hz. Delay and sound
// timers tick down at 60Hz, display refresh rate is 60Hz too. To get the 
// proper emulation speed and keep things simple, the main emulation loop runs
// at 60Hz. For each iteration we execute 540 / 60 = 9 instructions, then we
// decrease timers and refresh the display.
//
// The cycles variable is not part of the original CHIP-8 machine, it's just an
// artifact to keep track of how many instructions were executed in the current
// loop iteration.
type machine struct {
    breakpoints []uint16
    cycles      byte
    keyboard    [16]bool
    pixmap      [SCREENWIDTH][SCREENHEIGHT]uint8
    memory      [4096]byte
    regs        registers
    running     bool
    stack       [16]uint16
}

var m machine

var fonts [80]byte = [80]byte{
    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
    0x20, 0x60, 0x20, 0x20, 0x70, // 1
    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
    0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}


func machineAddBreakpoint(address uint16) {
    m.breakpoints = append(m.breakpoints, address)
}


func machineClearBreakpoints() {
    m.breakpoints = make([]uint16, 0, 5)
}


func machineDeleteBreakpoint(friendly int) {
    // friendly is the breakpoint id as seen by the user
    // id is the actual breakpoint id
    id := friendly - 1
    if friendly == 0 || id >= len(m.breakpoints) {
        fmt.Printf("Invalid breakpoint id\n")
    } else {
        for i := id; i < int(len(m.breakpoints)) - 1; i++ {
            m.breakpoints[i] = m.breakpoints[i + 1]
        }
        m.breakpoints = m.breakpoints[:len(m.breakpoints) - 1]
    }
}


func machineDisassembleInstruction(assembled uint16) (disassembled instruction) {
	// Extract variables from the assembled instruction
	// Obviously we won't need all of them, extra ones are just ignored
    nnn := assembled & 0x0FFF
    x := byte((assembled & 0x0F00) >> 8)
    y := byte((assembled & 0x00F0) >> 4)
    kk := byte(assembled & 0x00FF)
    n := byte(assembled & 0x000F)

    switch {
    case assembled == 0x00E0:
        return instruction{op:cls}

    case assembled == 0x00EE:
        return instruction{op:ret}

    case assembled & 0xF000 == 0x0000:
        return instruction{op:sys, nnn:nnn}

    case assembled & 0xF000 == 0x1000:
        return instruction{op:jmp, nnn:nnn}

    case assembled & 0xF000 == 0x2000:
        return instruction{op:call, nnn:nnn}

    case assembled & 0xF000 == 0x3000:
        return instruction{op:seb, x:x, kk:kk}

    case assembled & 0xF000 == 0x4000:
        return instruction{op:sneb, x:x, kk:kk}

    case assembled & 0xF000 == 0x5000:
        return instruction{op:ser, x:x, y:y}

    case assembled & 0xF000 == 0x6000:
        return instruction{op:ldb, x:x, kk:kk}

    case assembled & 0xF000 == 0x7000:
        return instruction{op:addb, x:x, kk:kk}

    case assembled & 0xF00F == 0x8000:
        return instruction{op:ldr, x:x, y:y}

    case assembled & 0xF00F == 0x8001:
        return instruction{op:or, x:x, y:y}

    case assembled & 0xF00F == 0x8002:
        return instruction{op:and, x:x, y:y}

    case assembled & 0xF00F == 0x8003:
        return instruction{op:xor, x:x, y:y}

    case assembled & 0xF00F == 0x8004:
        return instruction{op:addr, x:x, y:y}

    case assembled & 0xF00F == 0x8005:
        return instruction{op:sub, x:x, y:y}

    case assembled & 0xF00F == 0x8006:
        return instruction{op:shr, x:x, y:y}

    case assembled & 0xF00F == 0x8007:
        return instruction{op:subn, x:x, y:y}

    case assembled & 0xF00F == 0x800E:
        return instruction{op:shl, x:x, y:y}

    case assembled & 0xF000 == 0x9000:
        return instruction{op:sner, x:x, y:y}

    case assembled & 0xF000 == 0xA000:
        return instruction{op:ldi, nnn:nnn}

    case assembled & 0xF000 == 0xB000:
        return instruction{op:jpv, nnn:nnn}

    case assembled & 0xF000 == 0xC000:
        return instruction{op:rnd, x:x, kk:kk}

    case assembled & 0xF000 == 0xD000:
        return instruction{op:drw, x:x, y:y, n:n}

    case assembled & 0xF0FF == 0xE09E:
        return instruction{op:skp, x:x}

    case assembled & 0xF0FF == 0xE0A1:
        return instruction{op:sknp, x:x}

    case assembled & 0xF0FF == 0xF007:
        return instruction{op:gett, x:x}

    case assembled & 0xF0FF == 0xF00A:
        return instruction{op:ldk, x:x}

    case assembled & 0xF0FF == 0xF015:
        return instruction{op:sett, x:x}

    case assembled & 0xF0FF == 0xF018:
        return instruction{op:lds, x:x}

    case assembled & 0xF0FF == 0xF01E:
        return instruction{op:addi, x:x}

    case assembled & 0xF0FF == 0xF029:
        return instruction{op:ldf, x:x}

    case assembled & 0xF0FF == 0xF033:
        return instruction{op:ldbcd, x:x}

    case assembled & 0xF0FF == 0xF055:
        return instruction{op:save, x:x}

    case assembled & 0xF0FF == 0xF065:
        return instruction{op:restore, x:x}

    default:
		return instruction{}
        //TODO: return instruction, error instead of just instruction
		//TODO: return proper error in case instruction not recognized
    }
}


func machineGetInstruction(address uint16) uint16 {
    return uint16(m.memory[address]) << 8 + uint16(m.memory[address + 1])
}


func machineInitialize() {
    for i, v := range fonts {
        m.memory[MEMFONTS + i] = v
    }
    m.breakpoints = make([]uint16, 0, 5)
    machineReset()
}


func machineIsRunning() bool {
    return m.running
}


func machineLoadProgram(program string) {
    for i := MEMPROGRAMSTART; i < MEMEND; i++ {
        m.memory[i] = 0
    }

    data, err := ioutil.ReadFile(program)
    if err != nil {
	    log.Fatal(err)
    }

    // TODO: check if program is too big
    for i, v := range data {
        m.memory[MEMPROGRAMSTART + i] = v
    }
}


func machineListBreakpoints() []uint16{
    return m.breakpoints
}


func machinePlaySound() bool {
    if m.regs.st > 0 {
        return true
    }
    return false
}


func machineReset() {
    for i, _ := range m.keyboard {
        m.keyboard[i] = false
    }

	m.regs.i = 0
	m.regs.dt = 0
	m.regs.st = 0
	m.regs.pc = MEMPROGRAMSTART
    m.regs.sp = 16

    m.cycles = 0
    m.running = false
    rand.Seed(time.Now().UnixNano())

    for x := 0; x < 64; x++ {
        for y := 0; y < 32; y++ {
            m.pixmap[x][y] = 0
        }
    }
}


func machineRun(buzz chan struct{}, draw chan struct{}, stop chan struct{}) {
    m.running = true
    for {
        for i := m.cycles; i < 9; i++ {
            for i := 0; i < len(m.breakpoints); i++ {
                if m.regs.pc == m.breakpoints[i] {
                    fmt.Printf("Found breakpoint at 0x%03x\n", m.regs.pc)
                    m.running = false
                    return
                }
            }

            // we received a stop token from the CLI
            // dequeue it and exit the run
            if (len(stop)) == 1 {
                <-stop
                m.running = false
                return
            }

            machineStep(buzz, draw)
        }
        time.Sleep(SLEEPTIME)
    }
}


func machineStep(buzz chan struct{}, draw chan struct{}) {
    incrementPC := true
    instruction := machineDisassembleInstruction(machineGetInstruction(m.regs.pc))

    switch {
    case instruction.op == sys:
        // ignore and do nothing

    case instruction.op == cls:
        for x := 0; x < 64; x++ {
            for y := 0; y < 32; y++ {
                m.pixmap[x][y] = 0
            }
        }

    case instruction.op == ret:
        m.regs.pc = m.stack[m.regs.sp]
        m.stack[m.regs.sp] = 0 // Clean the value from m.stack, not required but better for debugging
        m.regs.sp++
        incrementPC = false

    case instruction.op == jmp:
        m.regs.pc = instruction.nnn
        incrementPC = false

    case instruction.op == call:
        // TODO: check if m.stack isn't full
        // TODO: check if addr is valid
        m.regs.sp--
        m.stack[m.regs.sp] = m.regs.pc + 2
        m.regs.pc = instruction.nnn
        incrementPC = false

    case instruction.op == seb:
        if m.regs.v[instruction.x] == instruction.kk {
            m.regs.pc += 2
        }

    case instruction.op == sneb:
        if m.regs.v[instruction.x] != instruction.kk {
            m.regs.pc += 2
        }

    case instruction.op == ser:
        if m.regs.v[instruction.x] == m.regs.v[instruction.y] {
            m.regs.pc += 2
        }

    case instruction.op == ldb:
        m.regs.v[instruction.x] = instruction.kk

    case instruction.op == addb:
        if int(m.regs.v[instruction.x]) + int(instruction.kk) > 256 {
            m.regs.v[0xf] = 1
        } else {
            m.regs.v[0xf] = 0
        }
        m.regs.v[instruction.x] += instruction.kk

    case instruction.op == ldr:
        m.regs.v[instruction.x] = m.regs.v[instruction.y]

    case instruction.op == or:
        m.regs.v[instruction.x] |= m.regs.v[instruction.y]

    case instruction.op == and:
        m.regs.v[instruction.x] &= m.regs.v[instruction.y]

    case instruction.op == xor:
        m.regs.v[instruction.x] ^= m.regs.v[instruction.y]

    case instruction.op == addr:
        m.regs.v[instruction.x] += m.regs.v[instruction.y]

    case instruction.op == sub:
        if m.regs.v[instruction.x] > m.regs.v[instruction.y] {
            m.regs.v[0xf] = 1
        } else {
            m.regs.v[0xf] = 0
        }
        m.regs.v[instruction.x] -= m.regs.v[instruction.y]

    case instruction.op == shr:
        if m.regs.v[instruction.x] & 0x01 == 1 {
            m.regs.v[0xf] = 1
        } else {
            m.regs.v[0xf] = 0
        }
        m.regs.v[instruction.x] = m.regs.v[instruction.x] >> 1

    case instruction.op == subn:
        if m.regs.v[instruction.y] > m.regs.v[instruction.x] {
            m.regs.v[0xf] = 1
        } else {
            m.regs.v[0xf] = 0
        }
        m.regs.v[instruction.x] = m.regs.v[instruction.y] - m.regs.v[instruction.x]

    case instruction.op == shl:
        if m.regs.v[instruction.x] >> 7 == 1 {
            m.regs.v[0xf] = 1
        } else {
            m.regs.v[0xf] = 0
        }
        m.regs.v[instruction.x] = m.regs.v[instruction.x] << 1

    case instruction.op == sner:
        if m.regs.v[instruction.x] != m.regs.v[instruction.y] {
            m.regs.pc += 2
        }

    case instruction.op == ldi:
        m.regs.i = instruction.nnn

    case instruction.op == jpv:
        m.regs.pc = uint16(m.regs.v[0x0]) + instruction.nnn
        incrementPC = false

    case instruction.op == rnd:
        m.regs.v[instruction.x] = byte(rand.Intn(255)) & instruction.kk

    case instruction.op == drw:
        m.regs.v[0xf] = 0

        for j := uint16(0); j < uint16(instruction.n); j++ {
            y := (uint16(m.regs.v[instruction.y]) + j)
            if y < SCREENHEIGHT {
                for i := uint16(0); i < 8; i++ {
                    x := (uint16(m.regs.v[instruction.x]) + i)
                    if x < SCREENWIDTH {
                        p := &m.pixmap[x][y]
                        n := (m.memory[m.regs.i + j] >> (8 - (i + 1))) & 0x1

                        old := *p
                        *p ^= n
                        if old > *p {
                            m.regs.v[0xf] = 1
                        }
                    }
                }
            }
        }

    case instruction.op == skp:
        if m.keyboard[m.regs.v[instruction.x]] {
            m.regs.pc += 2
        }

    case instruction.op == sknp:
        if ! m.keyboard[m.regs.v[instruction.x]] {
            m.regs.pc += 2
        }

    case instruction.op == gett:
        m.regs.v[instruction.x] = m.regs.dt

    case instruction.op == ldk:
        for {
            for i, _ := range m.keyboard {
                if m.keyboard[i] {
                    m.regs.v[instruction.x] = byte(i)
                    break
                }
            }
            time.Sleep(1000000)
        }

    case instruction.op == sett:
        m.regs.dt = m.regs.v[instruction.x]

    case instruction.op == lds:
        m.regs.st = m.regs.v[instruction.x]

    case instruction.op == addi:
        m.regs.i = m.regs.i + uint16(m.regs.v[instruction.x])

    case instruction.op == ldf:
        // TODO: check value in register is not bigger than 0xf
        m.regs.i = uint16(m.regs.v[instruction.x]) * 5

    case instruction.op == ldbcd:
        n := m.regs.v[instruction.x]
        m.memory[m.regs.i], n = n / 100, n % 100
        m.memory[m.regs.i + 1], n = n / 10, n % 10
        m.memory[m.regs.i + 2] = n

    case instruction.op == save:
        for j := uint16(0); j <= uint16(instruction.x); j++ {
            m.memory[m.regs.i + j] = m.regs.v[j]
        }

    case instruction.op == restore:
        for j := uint16(0); j <= uint16(instruction.x); j++ {
            m.regs.v[j] = m.memory[m.regs.i + j]
        }
    }

    if incrementPC {
        m.regs.pc += 2
    }

    m.cycles++
    if m.cycles == 9 {
	    if m.regs.dt > 0 {
		    m.regs.dt--
	    }
	    if m.regs.st > 0 {
	        m.regs.st--
	    }
        m.cycles = 0
        buzz <- struct{}{}
        draw <- struct{}{}
    }
}

func machineUpdateKeyboard(key byte, state bool) {
    m.keyboard[key] = state
}
