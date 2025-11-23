package main

import (
	"fmt"
	"strconv"
)

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type lights struct {
	ws wsEngine
}

func (lgt *lights) setup() error {
	return lgt.ws.Init()
}

func (lgt *lights) setFullStringColor(color uint32) error {
	for i := 0; i < len(lgt.ws.Leds(0)); i++ {
		lgt.ws.Leds(0)[i] = color
	}
	if err := lgt.ws.Render(); err != nil {
		return err
	}
	return nil
}

func (lgt *lights) setAltStringColor(color1, color2 uint32) error {
	for i := 0; i < len(lgt.ws.Leds(0)); i++ {
		if i%2 == 0 {
			lgt.ws.Leds(0)[i] = color1
		} else {
			lgt.ws.Leds(0)[i] = color2
		}
	}
	if err := lgt.ws.Render(); err != nil {
		return err
	}
	return nil
}

func hexToGRB(hexColor string) (uint32, error) {
	if len(hexColor) == 7 && hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}
	if len(hexColor) != 6 {
		return 0, fmt.Errorf("invalid hex color: %s", hexColor)
	}

	r, err := strconv.ParseInt(hexColor[0:2], 16, 32)
	if err != nil {
		return 0, err
	}
	g, err := strconv.ParseInt(hexColor[2:4], 16, 32)
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseInt(hexColor[4:6], 16, 32)
	if err != nil {
		return 0, err
	}

	grb := uint32(g)<<16 | uint32(r)<<8 | uint32(b)

	return grb, nil
}
