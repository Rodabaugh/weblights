package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Rodabaugh/weblights/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 90
)

type apiConfig struct {
	platform string
	db       *database.Queries
	lgts     *lights
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using enviroment variables.")
	} else {
		fmt.Println("Loaded .env file.")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	ledCount := os.Getenv("NUM_LEDS")
	if dbURL == "" {
		log.Fatal("NUM_LEDS must be set")
	}

	ledCountInt, err := strconv.ParseInt(ledCount, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	// Init lights
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = int(ledCountInt)

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
		db:   dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		MainPage(&apiCfg).Render(r.Context(), w)
	})
	mux.HandleFunc("GET /controls", func(w http.ResponseWriter, r *http.Request) {
		Controls(&apiCfg).Render(r.Context(), w)
	})
	mux.HandleFunc("GET /color-picker", func(w http.ResponseWriter, r *http.Request) {
		ColorPicker().Render(r.Context(), w)
	})

	mux.HandleFunc("POST /api/lights/color", apiCfg.setColor)
	mux.HandleFunc("POST /api/lights/altColor", apiCfg.setAltColor)
	mux.HandleFunc("POST /api/lights/preset", apiCfg.setPreset)
	mux.HandleFunc("POST /api/presets", apiCfg.newPreset)
	mux.HandleFunc("DELETE /api/presets", apiCfg.deletePreset)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting weblights in %s mode on port: 8080\n", platform)
	log.Fatal(server.ListenAndServe())
}
