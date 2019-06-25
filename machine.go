package main

// add package description

import (
    //"fmt"
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
	nnn uint16 // (or addr) A 12-bit value, the lowest 12 bits of the instruction<Paste>
	x byte     // A 4-bit value, the lower 4 bits of the high byte of the instruction
	y byte     // A 4-bit value, the upper 4 bits of the low byte of the instruction
	kk byte    // (or byte) An 8-bit value, the lowest 8 bits of the instruction
	n byte     // (or nibble) A 4-bit value, the lowest 4 bits of the instruction
}

type registers struct {
	v  [16]byte // Data registers V0 to VF
	i  uint16   // Address Register
	dt byte     // Delay Timer
	st byte     // Sound Timer
	pc uint16   // Program Counter
	sp byte     // Stack Pointer
}
var pixmap [SCREENWIDTH][SCREENHEIGHT]uint8
var memory [4096]byte
var stack [16]uint16
var state registers

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

func disassembleInstruction(assembled uint16) (disassembled instruction) {
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

func getInstruction(address uint16) uint16 {
    return uint16(memory[address]) << 8 + uint16(memory[address + 1])
}


func loadProgram(program string) {
    for i := MEMPROGRAMSTART; i < MEMEND; i++ {
        memory[i] = 0
    }

    data, err := ioutil.ReadFile(program)
    if err != nil {
	    log.Fatal(err)
    }

    // TODO: check if program is too big
    for i, v := range data {
        memory[MEMPROGRAMSTART + i] = v
    }
}

func runMachine() {
    for {
		// run loop at 120hz
		// every other iteration decrement timers
		// every iteration: calculate time ellapsed and compute number of instructions which should be run
        // So it's much better to run our loop at much lower speed (100Hz in my case), calculating time between two loop cycles, then based on target frequency calculate number of operations that should be performed and perform them at once.
        // When you are writing emulator, you have original CPU speed as a reference in most cases, but since CHIP-8 is interpreted language, speed varies based on device program was designed for. By my observations best universal speed for CHIP-8 programs is 500Hz and for SuperCHIP it's 1000hz, but you have to give user ability to change it so their experience is as good as possible. Also don't forget that delay and sound timers should always tick down at 60Hz, no matter how fast emulator is running.<Paste>


        stepMachine()
    }
}

func stepMachine() {
    incrementPC := true
    instruction := disassembleInstruction(getInstruction(state.pc))

    switch {
    case instruction.op == sys:
        // ignore and do nothing

    case instruction.op == cls:
        for x := 0; x < 64; x++ {
            for y := 0; y < 32; y++ {
                pixmap[x][y] = 0
            }
        }

    case instruction.op == ret:
        state.pc = stack[state.sp]
        stack[state.sp] = 0 // Clean the value from the stack, not required but better for debugging
        state.sp++

    case instruction.op == jmp:
        state.pc = instruction.nnn
        incrementPC = false

    case instruction.op == call:
        // TODO: check if stack isn't full
        // TODO: check if addr is valid
        state.sp--
        stack[state.sp] = state.pc + 2
        state.pc = instruction.nnn
        incrementPC = false

    case instruction.op == seb:
        if state.v[instruction.x] == instruction.kk {
            state.pc += 2
        }

    case instruction.op == sneb:
        if state.v[instruction.x] != instruction.kk {
            state.pc += 2
        }

    case instruction.op == ser:
        if state.v[instruction.x] == state.v[instruction.y] {
            state.pc += 2
        }

    case instruction.op == ldb:
        state.v[instruction.x] = instruction.kk

    case instruction.op == addb:
        if int(state.v[instruction.x]) + int(instruction.kk) > 256 {
            state.v[0xf] = 1
        } else {
            state.v[0xf] = 0
        }
        state.v[instruction.x] += instruction.kk

    case instruction.op == ldr:
        state.v[instruction.x] = state.v[instruction.y]

    case instruction.op == or:
        state.v[instruction.x] |= state.v[instruction.y]

    case instruction.op == and:
        state.v[instruction.x] &= state.v[instruction.y]

    case instruction.op == xor:
        state.v[instruction.x] ^= state.v[instruction.y]

    case instruction.op == addr:
        state.v[instruction.x] += state.v[instruction.y]

    case instruction.op == sub:
        if state.v[instruction.x] > state.v[instruction.y] {
            state.v[0xf] = 1
        } else {
            state.v[0xf] = 0
        }
        state.v[instruction.x] -= state.v[instruction.y]

    case instruction.op == shr:
        if state.v[instruction.x] & 0x01 == 1 {
            state.v[0xf] = 1
        } else {
            state.v[0xf] = 0
        }
        state.v[instruction.x] = state.v[instruction.x] >> 1

    case instruction.op == subn:
        if state.v[instruction.y] > state.v[instruction.x] {
            state.v[0xf] = 1
        } else {
            state.v[0xf] = 0
        }
        state.v[instruction.x] = state.v[instruction.y] - state.v[instruction.x]

    case instruction.op == shl:
        if state.v[instruction.x] >> 7 == 1 {
            state.v[0xf] = 1
        } else {
            state.v[0xf] = 0
        }
        state.v[instruction.x] = state.v[instruction.x] << 1

    case instruction.op == sner:
        if state.v[instruction.x] != state.v[instruction.y] {
            state.pc += 2
        }

    case instruction.op == ldi:
        state.i = instruction.nnn

    case instruction.op == jpv:
        state.pc = uint16(state.v[0x0]) + instruction.nnn
        incrementPC = false

    case instruction.op == rnd:
        state.v[instruction.x] = byte(rand.Intn(255)) & instruction.kk

    case instruction.op == drw:
        state.v[0xf] = 0

        for j := uint16(0); j < uint16(instruction.n); j++ {
            y := (uint16(state.v[instruction.y]) + j)
            if y < SCREENHEIGHT {
                for i := uint16(0); i < 8; i++ {
                    x := (uint16(state.v[instruction.x]) + i)
                    if x < SCREENWIDTH {
                        p := &pixmap[x][y]
                        n := (memory[state.i + j] >> (8 - (i + 1))) & 0x1

                        old := *p
                        *p ^= n
                        if old > *p {
                            state.v[0xf] = 1
                        }
                    }
                }
            }
        }

    case instruction.op == skp:
        // TODO: implement me!

    case instruction.op == sknp:
        // TODO: implement me!

    case instruction.op == gett:
        state.v[instruction.x] = state.dt

    case instruction.op == ldk:
        // TODO: implement me!

    case instruction.op == sett:
        state.dt = state.v[instruction.x]

    case instruction.op == lds:
        state.st = instruction.x

    case instruction.op == addi:
        state.i = state.i + uint16(state.v[instruction.x])

    case instruction.op == ldf:
        // TODO: check value in register is not bigger than 0xf
        state.i = uint16(state.v[instruction.x]) * 5

    case instruction.op == ldbcd:
        n := state.v[instruction.x]
        memory[state.i], n = n / 100, n % 100
        memory[state.i + 1], n = n / 10, n % 10
        memory[state.i + 2] = n

    case instruction.op == save:
        for j := uint16(0); j < uint16(instruction.x); j++ {
            memory[state.i + j] = state.v[j]
        }

    case instruction.op == restore:
        for j := uint16(0); j < uint16(instruction.x); j++ {
            state.v[j] = memory[state.i + j]
        }
    }

    if incrementPC {
        state.pc += 2
    }

	//TODO: decrement timers for debugging purposes only
	if state.dt > 0 {
		state.dt--
	}
	if state.st > 0 {
		state.st--
	}
}

func initializeMachine() {
    resetRegisters()

    for i, v := range fonts {
        memory[MEMFONTS + i] = v
    }

    rand.Seed(time.Now().UnixNano())
}

func resetRegisters() {
    for i, _ := range state.v {
        state.v[i] = 0
    }
	state.i = 0
	state.dt = 0
	state.st = 0
	state.pc = MEMPROGRAMSTART
    state.sp = 16
}
