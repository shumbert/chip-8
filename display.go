package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "time"
)

const (
    MAGNIFICATION = 8
    SCREENWIDTH = 64
    SCREENHEIGHT = 32
)

var display [SCREENWIDTH][SCREENHEIGHT]uint8

var window *sdl.Window
var surface *sdl.Surface
var black, white uint32

func cleanupDisplay() {
    window.Destroy()
    sdl.Quit()
}

func initializeDisplay() {
    // TODO: assign err and process it
    sdl.Init(sdl.INIT_EVERYTHING)
    window, _ = sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, SCREENWIDTH * MAGNIFICATION, SCREENHEIGHT * MAGNIFICATION, sdl.WINDOW_SHOWN)
    surface, _ = window.GetSurface()
    black = sdl.MapRGB(surface.Format, 0x00, 0x00, 0x00)
    white = sdl.MapRGB(surface.Format, 0xff, 0xff, 0xff)
}

func redrawDisplay() {
    surface.FillRect(nil, black)
    window.UpdateSurface()
}

func runDisplay() {
    initializeDisplay()
    defer cleanupDisplay()

    redrawDisplay()

    for {
        time.Sleep(1 * time.Second)
    }
}
