package main

// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"
import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
    "log"
    "math"
    "reflect"
    "unsafe"
)

const (
    MAGNIFICATION = 8
	SAMPLEHZ      = 48000
    TONEHZ        = 440
	DPHASE        = 2 * math.Pi * TONEHZ / SAMPLEHZ
)

var window *sdl.Window
var surface, pixelSurface *sdl.Surface
var pixelRect *sdl.Rect
var black, white uint32
var sample_nr int

//export SineWave
func SineWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	var phase float64
	for i := 0; i < n; i += 2 {
		phase += DPHASE
		sample := C.Uint8((math.Sin(phase) + 0.999999) * 128)
		buf[i] = sample
		buf[i+1] = sample
	}
}

func ioCleanupDisplay() {
    window.Destroy()
    sdl.CloseAudio()
    sdl.Quit()
}


func ioInit() {
    // TODO: assign err and process it
    sdl.Init(sdl.INIT_EVERYTHING)

    window, _ = sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, SCREENWIDTH * MAGNIFICATION, SCREENHEIGHT * MAGNIFICATION, sdl.WINDOW_SHOWN)
    surface, _ = window.GetSurface()

    black = sdl.MapRGB(surface.Format, 0x00, 0x00, 0x00)
    white = sdl.MapRGB(surface.Format, 0xff, 0xff, 0xff)

    pixelSurface, _ = sdl.CreateRGBSurface(0, MAGNIFICATION, MAGNIFICATION, 8, 0, 0, 0, 0)
    pixelSurface.FillRect(nil, white)

    pixelRect = new(sdl.Rect)
    pixelRect.W = MAGNIFICATION
    pixelRect.H = MAGNIFICATION

    spec := &sdl.AudioSpec{
		Freq:     SAMPLEHZ,
		Format:   sdl.AUDIO_U8,
		Channels: 2,
		Samples:  SAMPLEHZ,
		Callback: sdl.AudioCallback(C.SineWave),
	}

	if err := sdl.OpenAudio(spec, nil); err != nil {
		log.Println(err)
		return
	}
}


func ioRedrawDisplay() {
    surface.FillRect(nil, black)

    for x := 0; x < SCREENWIDTH; x++ {
        for y := 0; y < SCREENHEIGHT; y++ {
            if m.pixmap[x][y] == 1 {
                pixelRect.X = int32(x * MAGNIFICATION)
                pixelRect.Y = int32(y * MAGNIFICATION)
                pixelSurface.Blit(nil, surface, pixelRect)
            }
        }
    }
    window.UpdateSurface()
}


func ioRunBuzzer(buzz chan struct{}) {
    for {
        <-buzz
        fmt.Printf("buzzer is alive: %d %t\n", m.regs.st, m.running)
        if machinePlaySound() && machineIsRunning() {
            sdl.PauseAudio(false)
        } else {
            sdl.PauseAudio(true)
        }
    }
}

func ioRunDisplay(draw chan struct{}) {
    defer ioCleanupDisplay()

    ioRedrawDisplay()
    for {
        <-draw
        ioRedrawDisplay()
    }
}

func ioRunKeyboard() {
    var e sdl.Event
    var k byte

    for {
        e = sdl.WaitEvent()
        if e != nil {

            switch e.(type) {
		    case *sdl.KeyboardEvent:

                // We are not interested in repeat events
				if e.(*sdl.KeyboardEvent).Repeat > 0 {
                    continue
                }

				switch e.(*sdl.KeyboardEvent).Keysym.Sym {
				case sdl.K_1:
                    k = 1        // 1 maps to 1
				case sdl.K_2:
                    k = 2        // 2 maps to 2
				case sdl.K_3:
                    k = 3        // 3 maps to 3
				case sdl.K_4:
                    k = 12       // 4 maps to C
				case sdl.K_q:
                    k = 4        // q maps to 4
				case sdl.K_w:
                    k = 5        // w maps to 5
				case sdl.K_e:
                    k = 6        // e maps to 6
				case sdl.K_r:
                    k = 13       // r maps to D
				case sdl.K_a:
                    k = 7        // a maps to 7
				case sdl.K_s:
                    k = 8        // s maps to 8
				case sdl.K_d:
                    k = 9        // d maps to 9
				case sdl.K_f:
                    k = 14       // f maps to E
				case sdl.K_z:
                    k = 10       // z maps to A
				case sdl.K_x:
                    k = 0        // x maps to 0
				case sdl.K_c:
                    k = 11       // c maps to B
				case sdl.K_v:
                    k = 15       // v maps to F
                }

				switch e.(*sdl.KeyboardEvent).State {
                case sdl.PRESSED:
                    machineUpdateKeyboard(k, true)
                case sdl.RELEASED:
                    machineUpdateKeyboard(k, false)
                }

            }

        }
    }
}
