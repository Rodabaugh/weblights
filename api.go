package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rodabaugh/weblights/internal/database"
	"github.com/google/uuid"
)

type Preset struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Colors    []int64
}

func (apiCfg apiConfig) setColor(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Color string `json:"color"`
	}

	type response struct {
		success bool `json:"success"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Solid Color: %s", params.Color), false)
		return
	}

	newColor, err := hexToGRB(params.Color)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Solid Color: %s", params.Color), false)
		Fail().Render(r.Context(), w)
		return
	}

	fmt.Printf("Settings lights to %s\n", params.Color)

	err = apiCfg.lgts.setFullStringColor(newColor)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Solid Color: %s", params.Color), false)
		Fail().Render(r.Context(), w)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			success: true,
		})
		return
	}

	apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Solid Color: %s", params.Color), true)
	Success().Render(r.Context(), w)
}

func (apiCfg apiConfig) setAltColor(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Color1 string `json:"color1"`
		Color2 string `json:"color2"`
	}

	type response struct {
		success bool `json:"success"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Alt Colors: %s, %s", params.Color1, params.Color2), false)
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	var colors []int64

	newColor1, err := hexToGRB(params.Color1)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Alt Colors: %s, %s", params.Color1, params.Color2), false)
		Fail().Render(r.Context(), w)
		return
	}

	newColor2, err := hexToGRB(params.Color2)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Alt Colors: %s, %s", params.Color1, params.Color2), false)
		Fail().Render(r.Context(), w)
		return
	}

	colors = append(colors, newColor1, newColor2)

	fmt.Printf("Settings lights to %s and %s\n", params.Color1, params.Color2)

	err = apiCfg.lgts.setAltStringColor(colors)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Alt Colors: %s, %s", params.Color1, params.Color2), false)
		Fail().Render(r.Context(), w)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			success: true,
		})
		return
	}

	apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Alt Colors: %s, %s", params.Color1, params.Color2), true)
	Success().Render(r.Context(), w)
}

func (apiCfg apiConfig) getPresets() []Preset {
	databasePresets, err := apiCfg.db.GetAllPresets(context.Background())
	if err != nil {
		fmt.Printf("unable to get presets from database", err)
		return []Preset{}
	}

	presets := []Preset{}

	for _, dbPreset := range databasePresets {
		presets = append(presets, Preset{
			ID:        dbPreset.ID,
			CreatedAt: dbPreset.CreatedAt,
			UpdatedAt: dbPreset.UpdatedAt,
			Name:      dbPreset.Name,
			Colors:    dbPreset.Colors,
		})
	}

	return presets
}

func (apiCfg apiConfig) newPreset(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name   string   `json:"name"`
		Colors []string `json:"colors"`
	}

	type response struct {
		success bool `json:"success"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s", params.Name), false)
		return
	}

	var colors []int64

	for _, color := range params.Colors {
		colorInt, err := hexToGRB(color)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
			apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s", params.Name), false)
			return
		}
		colors = append(colors, colorInt)
	}

	newDBPreset, err := apiCfg.db.CreatePreset(r.Context(), database.CreatePresetParams{
		Name:   fmt.Sprintf("%s - %+v", params.Name, params.Colors),
		Colors: colors,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was write to DB", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s", params.Name), false)
		return
	}

	apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s", params.Name), true)
	SetStatus(fmt.Sprintf("Created new preset: %s", newDBPreset.Name)).Render(r.Context(), w)
}

func (apiCfg apiConfig) setPreset(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		PresetID string `json:"presetId"`
	}

	type response struct {
		success bool `json:"success"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Set Preset: %s", params.PresetID), false)
		return
	}

	uuid, err := uuid.Parse(params.PresetID)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Set Preset: Invalid UUID %s", params.PresetID), false)
		respondWithError(w, http.StatusBadRequest, "Invalid preset ID (not a valid UUID)", err)
		return
	}

	fmt.Println(params.PresetID)
	preset, err := apiCfg.db.GetPresetByID(r.Context(), uuid)

	fmt.Printf("Settings lights to %s\n", preset.Name)

	err = apiCfg.lgts.setAltStringColor(preset.Colors)
	if err != nil {
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Set Preset: %s", params.PresetID), false)
		Fail().Render(r.Context(), w)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			success: true,
		})
		return
	}

	apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("Set Preset: %s", params.PresetID), true)
	SetStatus(fmt.Sprintf("Successfully set preset to: %s", preset.Name)).Render(r.Context(), w)
}

func (apiCfg apiConfig) newLogEntry(ctx context.Context, requester, request string, result bool) {
	_, err := apiCfg.db.CreateLogEntry(ctx, database.CreateLogEntryParams{
		Requester: requester,
		Request:   request,
		Result:    result,
	})
	if err != nil {
		fmt.Printf("DB Error: %v\n", err)
	}
}
