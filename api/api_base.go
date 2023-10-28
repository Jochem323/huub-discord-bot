package api

import (
	"encoding/json"
	"fmt"
	"huub-discord-bot/common"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api/guilds", MakeHandlerFunc(s.HandleGetGuilds)).Methods(http.MethodGet)
	router.HandleFunc("/api/guilds", MakeHandlerFunc(s.HandleCreateGuild)).Methods(http.MethodPost)
	router.HandleFunc("/api/guilds/{guild_id}", MakeHandlerFunc(s.HandleGetGuild)).Methods(http.MethodGet)
	router.HandleFunc("/api/guilds/{guild_id}", MakeHandlerFunc(s.HandleUpdateGuild)).Methods(http.MethodPut)
	router.HandleFunc("/api/guilds/{guild_id}", MakeHandlerFunc(s.HandleDeleteGuild)).Methods(http.MethodDelete)

	router.HandleFunc("/api/keywords", MakeHandlerFunc(s.HandleGetKeywords)).Methods(http.MethodGet)
	router.HandleFunc("/api/keywords", MakeHandlerFunc(s.HandleCreateKeyword)).Methods(http.MethodPost)
	router.HandleFunc("/api/keywords/{keyword_id}", MakeHandlerFunc(s.HandleGetKeyword)).Methods(http.MethodGet)
	router.HandleFunc("/api/keywords/{keyword_id}", MakeHandlerFunc(s.HandleUpdateKeyword)).Methods(http.MethodPut)
	router.HandleFunc("/api/keywords/{keyword_id}", MakeHandlerFunc(s.HandleDeleteKeyword)).Methods(http.MethodDelete)

	router.HandleFunc("/api/apikeys", MakeHandlerFunc(s.HandleGetKeys)).Methods(http.MethodGet)
	router.HandleFunc("/api/apikeys", MakeHandlerFunc(s.HandleCreateKey)).Methods(http.MethodPost)
	router.HandleFunc("/api/apikeys/{key_id}", MakeHandlerFunc(s.HandleGetKey)).Methods(http.MethodGet)
	router.HandleFunc("/api/apikeys/{key_id}", MakeHandlerFunc(s.HandleUpdateKey)).Methods(http.MethodPut)
	router.HandleFunc("/api/apikeys/{key_id}", MakeHandlerFunc(s.HandleDeleteKey)).Methods(http.MethodDelete)

	log.Println("API Server running")

	http.ListenAndServe(s.ListenAdress, router)
}

func MakeHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			SendMessageResponse(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func SendEmptyResponse(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}

func SendMessageResponse(w http.ResponseWriter, status int, msg string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(APIResponse{Message: msg})
}

func SendDataResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) GetKeyFromRequest(r *http.Request) (common.APIKey, error) {
	secret := os.Getenv("JWT_SECRET")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return common.APIKey{}, fmt.Errorf("no token provided")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return common.APIKey{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	id := int(claims["id"].(float64))

	key, err := s.APIKeyStore.GetKey(id)
	if err != nil {
		return common.APIKey{}, err
	}

	if key.Revoked {
		key.Active = false
		s.APIKeyStore.UpdateKey(key)
	}

	return key, nil
}

func (s *APIServer) GetKeyValidity(key common.APIKey) bool {
	if key.Admin {
		return key.Active
	}

	guild, err := s.GuildStore.GetGuild(key.GuildID)
	if err != nil {
		return false
	}

	if !guild.APIEnabled {
		return false
	}

	return key.Active
}
