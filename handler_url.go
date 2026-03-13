package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mosgizy/url-shortner/internal/database"
)
func (apiCfg *apiConfig) handlerCreateUrl(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Url string `json:"url"`
		// Code string `json:"code"` 
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Url == "" {
		respondWithError(w,400,"URL is required")
		return
	}

	shortCode := GenerateSortCode(6)

	// if shortCode == "" {
	// 	shortCode = GenerateSortCode(6)
	// }

	url, err := apiCfg.DB.CreateUrl(r.Context(),database.CreateUrlParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ShortCode: shortCode,
		LongUrl: params.Url,
		Clicks: 0,
	})
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Error creating URL: %v", err))
		return
	}

	respondWithJSON(w,201,url)
}	

func (apiCfg *apiConfig) handlerRedirectUrl(w http.ResponseWriter, r *http.Request){
	shortCode := chi.URLParam(r,"shortCode")

	url, err := apiCfg.DB.GetByShortCode(r.Context(),shortCode)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("URL not found: %v", err))
		return
	}

	err = apiCfg.DB.IncrementClicks(r.Context(),shortCode)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Unable to increment clicks: %v", err))
		return
	}

	http.Redirect(w,r,url.LongUrl,http.StatusMovedPermanently)
}