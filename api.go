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
	Color1    uint32
	Color2    uint32
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

	fmt.Printf("Settings lights to %s and %s\n", params.Color1, params.Color2)

	err = apiCfg.lgts.setAltStringColor(newColor1, newColor2)
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
			Color1:    uint32(dbPreset.Color1),
			Color2:    uint32(dbPreset.Color2),
		})
	}

	return presets
}

func (apiCfg apiConfig) newPreset(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name   string `json:"name"`
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
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s - %s, %s", params.Name, params.Color1, params.Color2), false)
		return
	}

	newColor1, err := hexToGRB(params.Color1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s - %s, %s", params.Name, params.Color1, params.Color2), false)
		return
	}

	newColor2, err := hexToGRB(params.Color2)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s - %s, %s", params.Name, params.Color1, params.Color2), false)
		return
	}

	newDBPreset, err := apiCfg.db.CreatePreset(r.Context(), database.CreatePresetParams{
		Name:   fmt.Sprintf("%s - %s - %s", params.Name, params.Color1, params.Color2),
		Color1: int64(newColor1),
		Color2: int64(newColor2),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was write to DB", err)
		apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s - %s, %s", params.Name, params.Color1, params.Color2), false)
		return
	}

	apiCfg.newLogEntry(r.Context(), r.RemoteAddr, fmt.Sprintf("New Preset: %s - %s, %s", params.Name, params.Color1, params.Color2), true)
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

	err = apiCfg.lgts.setAltStringColor(uint32(preset.Color1), uint32(preset.Color2))
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
