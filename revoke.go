package main

import (
	"http_adv/internal/auth"
	"net/http"
)

func (cfg *apiConfig)handlerRevoke(w http.ResponseWriter, r *http.Request){
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}
	err = cfg.database.RevokeRefreshToken(r.Context(),refreshToken)
	if err!=nil{
		//give an error message
		return
	}
	w.WriteHeader(http.StatusNoContent)
}