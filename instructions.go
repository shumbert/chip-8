package main

import (
    "fmt"
)

// Fx1E - ADD I, Vx
func insAddMemReg(x uint16) {
    state.i = state.i + x
}

// 8xy4 - ADD Vx, Vy
func insAddReg(x, y uint16) {
    state.v[x] += state.v[y]
}

// 7xkk - ADD Vx, byte
func insAddVal(x, val uint16) {
    if int(state.v[x]) + int(val) > 256 {
        state.v[0xf] = 1
    } else {
        state.v[0xf] = 0
    }
    state.v[x] += val
}

// 8xy2 - AND Vx, Vy
func insAnd(x, y uint16) {
    state.v[x] &= state.v[y]
}

// 2nnn - CALL addr
func insCall(addr uint16) {
    // TODO: check if stack isn't full
    // TODO: check if addr is valid
    state.sp--
    stack[state.sp] = state.pc
    state.pc = addr
}

// 00E0 - CLS
func insClearDisplay() {
    for x := 0; x < 64; x++ {
        for y := 0; y < 32; y++ {
            display[x][y] = 0
        }
    }
}

func insDisplaySprite() {
    // TODO: implement me!!!
    fmt.Println("insDisplaySprite called")
}

// Fx07 - LD Vx, DT
func insGetDelay(x uint16) {
    state.v[x] = state.delay
}

// 1nnn - JP addr
func insJump(addr uint16) {
    state.pc = addr
}

// Bnnn - JP V0, addr
func insJumpReg(addr uint16) {
    state.pc = state.v[0x0] + addr
}

// Annn - LD I, addr
func insLoadMemReg(addr uint16) {
    state.i = addr
}

// 8xy0 - LD Vx, Vy
func insLoadReg(x, y uint16) {
    state.v[x] = state.v[y]
}

// 6xkk - LD Vx, byte
func insLoadVal(x, val uint16) {
    state.v[x] = val
}

func insNope() {
    fmt.Println("insNope called")
}

// 8xy1 - OR Vx, Vy
func insOr(x, y uint16) {
    state.v[x] |= state.v[y]
}

// Cxkk - RND Vx, byte
func insRandom(x, val uint16) {
    state.v[x] = rand.Intn(255) & val
}

// Fx65 - LD Vx, [I]
func insRestoreRegisters(x uint16) {
    for j := 0; j < x; j++ {
        state.v[j] = memory[state.i + j]
    }
}

// 00EE - RET
func insReturn() {
    state.pc = stack[state.sp]
    state.sp++
}

// Fx55 - LD [I], Vx
func insSaveRegisters(x uint16) {
    for j := 0; j < x; j++ {
        memory[state.i + j] = state.v[j]
    }
}

// Fx15 - LD DT, Vx
func insSetDelay(x uint16) {
    state.delay = x
}

// Fx18 - LD ST, Vx
func insSetSound(x uint16) {
    state.sound = x
}

// 8xyE - SHL Vx {, Vy}
func insShiftLeft() {
    fmt.Println("insShiftLeft called")
}

func insShiftRight() {
    fmt.Println("insShiftRight called")
}

func insSkipEqualReg() {
    fmt.Println("insSkipEqualReg called")
}

func insSkipEqualVal() {
    fmt.Println("insSkipEqualVal called")
}

func insSkipKeyNotPressed() {
    fmt.Println("insSkipKeyNotPressed called")
}

func insSkipKeyPressed() {
    fmt.Println("insSkipKeyPressed called")
}

func insSkipNonEqualReg() {
    fmt.Println("insSkipNonEqualReg called")
}

func insSkipNonEqualVal() {
    fmt.Println("insSkipNonEqualVal called")
}

func insSpriteToMemReg() {
    fmt.Println("insSpriteToMemReg called")
}

func insStoreBCD() {
    fmt.Println("insStoreBCD called")
}

func insSub() {
    fmt.Println("insSub called")
}

func insSubNegate() {
    fmt.Println("insSubNegate called")
}

func insSys() {
    fmt.Println("insSys called")
}

func insWaitKey() {
    fmt.Println("insWaitKey called")
}

func insXor() {
    fmt.Println("insXor called")
}
