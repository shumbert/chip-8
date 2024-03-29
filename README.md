CHIP-8 emulator written in Golang with libsdl.

# Current status
- pong works
- untested with other roms
- sound is buggy/crappy
 
## todo
- try out with other roms
- fix sound
- add proper error handling where needed
- fix packaging and try on other platforms
- implement cli command history
- improve disassembly by parsing the binary to check which memory locations are used for sprites

# Resources
## Go
- [The Go wiki](https://github.com/golang/go/wiki)
- [Awesome Go](https://awesome-go.com/)
- [The Go Tour](http://tour.golang.org/)
- [How to write Go code](https://golang.org/doc/code.html)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Packages](https://golang.org/pkg/)
- [Writing a simple shell in Go](https://sj14.gitlab.io/post/2018-07-01-go-unix-shell/)

## LibSDL
- [Beginning Game Programming v2.0](http://lazyfoo.net/tutorials/SDL/index.php)
- [GoDoc - Package SDL](https://godoc.org/github.com/veandco/go-sdl2/sdl)
- [Another app built with golang and libsdl](https://github.com/fiorix/cat-o-licious)

# Emulator
- [BUILDING 8-BIT EMULATOR IN GOLANG](https://engineering.wpengine.com/building-8-bit-emulator-in-golang/)
- [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)
- [BYTE magazine, December 1978, An Easy Programming System](https://archive.org/details/byte-magazine-1978-12/page/n109)
- [How to write an emulator (CHIP-8 interpreter)](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
- [Cowgod's Chip-8 Technical Reference v1.0](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [Mastering Chip-8](http://mattmik.com/files/chip8/mastering/chip8.html)
- [Assembler and Emulator in Go](https://massung.github.io/CHIP-8/)
- [CHIP-8 emulatoar](https://colineberhardt.github.io/wasm-rust-chip8/web/)
- [Octo is a high level assembler, disassembler and simulator for Chip8](https://johnearnest.github.io/Octo/)

# ROMs
You can get ROMs [here](https://github.com/dmatlack/chip8/tree/master/roms).

# Requirements
Install libsdl2:
```
sudo apt install libsdl2-dev
go get -v github.com/veandco/go-sdl2/sdl
```

# Run the emulator
```
go run cli.go io.go machine.go main.go ~/Documents/Geek/Projects/go/Pong\ \[Paul\ Vervalin\,\ 1990\].ch8
```

# Emulation speed
https://www.reddit.com/r/EmuDev/comments/9hx3ry/how_to_do_timing/

When you are writing emulator, you have original CPU speed as a reference in most cases, but since CHIP-8 is interpreted language, speed varies based on device program was designed for. By my observations best universal speed for CHIP-8 programs is 500Hz and for SuperCHIP it's 1000hz, but you have to give user ability to change it so their experience is as good as possible. Also don't forget that delay and sound timers should always tick down at 60Hz, no matter how fast emulator is running.

# Octo
Octo provides basic debugging facilities for Chip8 programs. While a program is running, pressing the “i” key will interrupt execution and display the contents of the v registers, i and the program counter. Any register aliases and (guessed) labels will be indicated next to the raw register contents. You can click on registers in this view to cycle through displaying their contents in binary, decimal, or hexadecimal.

When interrupted, pressing “i” again or clicking the “continue” icon will resume execution, while pressing “o” will single-step through the program. The “u” key will attempt to step out (execute until the current subroutine returns) and the “l” key will attempt to step over (execute the contents of any subroutines until they return to the current level).

Breakpoints can also be placed in source code by using the command :breakpoint followed by a name- the name will be shown when the breakpoint is encountered so that multiple breakpoints can be readily distinguished. :breakpoint is an out-of-band debugging facility and inserting a breakpoint into your program will not add any code or modify any Chip8 registers.

# Keyboard
The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad with the following layout:
```
1 2 3 C
4 5 6 D
7 8 9 E
A 0 B F
```

This layout is mapped to the following keys (qwerty keyboard):
```
1 2 3 4
q w e r
a s d f
z x c v
```
