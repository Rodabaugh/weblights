package main

import (
	"log"
	"net/http"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

// Replace with env vars later
const (
	brightness = 90
	ledCounts  = 250
)

type apiConfig struct {
	platform string
	lgts     *lights
}

func main() {
	// Init lights
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	lights := &lights{
		ws: dev,
	}

	checkError(lights.setup())
	defer dev.Fini()

	// Init server
	apiCfg := apiConfig{
		lgts: lights,
	}

	// Testing
	lights.setFullStringColor(uint32(0x6feb92))
	time.Sleep(time.Millisecond * 100)
	lights.setFullStringColor(uint32(0xc1f677))
	time.Sleep(time.Millisecond * 100)
	lights.setFullStringColor(uint32(0x74318f))
	time.Sleep(time.Millisecond * 100)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		MainPage().Render(r.Context(), w)
	})

	mux.HandleFunc("POST /api/color", apiCfg.setColor)
	mux.HandleFunc("POST /api/altColor", apiCfg.setAltColor)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting weblights on port: 8080\n")
	log.Fatal(server.ListenAndServe())
}
