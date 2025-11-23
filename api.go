package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		return
	}

	newColor, err := hexToGRB(params.Color)
	if err != nil {
		Fail().Render(r.Context(), w)
		return
	}

	fmt.Printf("Settings lights to %s\n", params.Color)

	err = apiCfg.lgts.setFullStringColor(newColor)
	if err != nil {
		Fail().Render(r.Context(), w)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			success: true,
		})
		return
	}
	Success().Render(r.Context(), w)
}
