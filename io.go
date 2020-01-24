package main

import (
    "github.com/veandco/go-sdl2/sdl"
)

const (
    MAGNIFICATION = 8
)

var window *sdl.Window
var surface, pixelSurface *sdl.Surface
var pixelRect *sdl.Rect
var black, white uint32


func ioCleanupDisplay() {
    window.Destroy()
    sdl.Quit()
}


func ioInitDisplay() {
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


func ioRunDisplay(draw chan struct{}) {
    ioInitDisplay()
    defer ioCleanupDisplay()

    ioRedrawDisplay()
    for {
        <-draw
        ioRedrawDisplay()
    }
}
