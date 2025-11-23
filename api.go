package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rodabaugh/weblights/internal/database"
)

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
