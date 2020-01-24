# Resources
## Go
- [The Go wiki] (https://github.com/golang/go/wiki)
- [Awesome Go](https://awesome-go.com/)
- [The Go Tour](http://tour.golang.org/)
- [How to write Go code](https://golang.org/doc/code.html)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Packages](https://golang.org/pkg/)
- https://sj14.gitlab.io/post/2018-07-01-go-unix-shell/

## LibSDL
- http://lazyfoo.net/tutorials/SDL/index.php
- https://godoc.org/github.com/veandco/go-sdl2/sdl

# Environment
Add $GOPATH/bin to the PATH:
```
PATH="$PATH:$(go env GOPATH)/bin"
```

# Requirements
Install libsdl2:
```
sudo apt install libsdl2-dev
go get -v github.com/veandco/go-sdl2/sdl
```

# Emulator
- [BUILDING 8-BIT EMULATOR IN GOLANG](https://engineering.wpengine.com/building-8-bit-emulator-in-golang/)
- [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)
- [BYTE magazine, December 1978, An Easy Programming System](https://archive.org/details/byte-magazine-1978-12/page/n109)
- [How to write an emulator (CHIP-8 interpreter)](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
- https://github.com/dmatlack/chip8/tree/master/roms
- [Cowgod's Chip-8 Technical Reference v1.0](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [Mastering Chip-8](http://mattmik.com/files/chip8/mastering/chip8.html)
- [Assembler and Emulator in Go](https://massung.github.io/CHIP-8/)
- [](https://colineberhardt.github.io/wasm-rust-chip8/web/)
- [](https://johnearnest.github.io/Octo/)


# Run the emulator

# Emulation speed
https://www.reddit.com/r/EmuDev/comments/9hx3ry/how_to_do_timing/

When you are writing emulator, you have original CPU speed as a reference in most cases, but since CHIP-8 is interpreted language, speed varies based on device program was designed for. By my observations best universal speed for CHIP-8 programs is 500Hz and for SuperCHIP it's 1000hz, but you have to give user ability to change it so their experience is as good as possible. Also don't forget that delay and sound timers should always tick down at 60Hz, no matter how fast emulator is running.

# Octo
Octo provides basic debugging facilities for Chip8 programs. While a program is running, pressing the “i” key will interrupt execution and display the contents of the v registers, i and the program counter. Any register aliases and (guessed) labels will be indicated next to the raw register contents. You can click on registers in this view to cycle through displaying their contents in binary, decimal, or hexadecimal.

When interrupted, pressing “i” again or clicking the “continue” icon will resume execution, while pressing “o” will single-step through the program. The “u” key will attempt to step out (execute until the current subroutine returns) and the “l” key will attempt to step over (execute the contents of any subroutines until they return to the current level).

Breakpoints can also be placed in source code by using the command :breakpoint followed by a name- the name will be shown when the breakpoint is encountered so that multiple breakpoints can be readily distinguished. :breakpoint is an out-of-band debugging facility and inserting a breakpoint into your program will not add any code or modify any Chip8 registers.

# TODO
- rename all functions so that the name identifies which file they're in, i.e. runMachine() to machineRun()
- add exceptions where needed:
  - typically machineDeleteBreakpoint should return an error if the breakpoint id is not valid
