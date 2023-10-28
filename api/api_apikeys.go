package api

import (
	"encoding/json"
	"huub-discord-bot/common"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func (s *APIServer) HandleGetKeys(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	if !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	APIkeys, err := s.APIKeyStore.GetKeys()
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, APIkeys)
}

func (s *APIServer) HandleGetKey(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	if !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedKey, err := strconv.Atoi(mux.Vars(r)["key_id"])
	if err != nil {
		return err
	}

	APIKey, err := s.APIKeyStore.GetKey(requestedKey)
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, APIKey)
}

func (s *APIServer) HandleCreateKey(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	if !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	body := new(CreateKeyRequest)
	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return err
	}

	APIKey := common.APIKey{
		Admin:     body.Admin,
		GuildID:   body.GuildID,
		Comment:   body.Comment,
		CreatedBy: body.CreatedBy,
		CreatedAt: time.Now(),
		Active:    true,
		Revoked:   false,
		Ratelimit: true,
	}

	id, err := s.APIKeyStore.AddKey(APIKey)
	if err != nil {
		return err
	}

	APIKey.ID = id

	token, err := GenerateJWT(&APIKey)
	if err != nil {
		return err
	}

	response := CreateKeyResponse{
		Key:   APIKey,
		Token: token,
	}

	return SendDataResponse(w, http.StatusCreated, response)
}

func (s *APIServer) HandleUpdateKey(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	if !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	body := new(UpdateKeyRequest)
	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return err
	}

	requestedKey, err := strconv.Atoi(mux.Vars(r)["key_id"])
	if err != nil {
		return err
	}

	APIKey, err := s.APIKeyStore.GetKey(requestedKey)
	if err != nil {
		return err
	}

	if body.Comment != "" {
		APIKey.Comment = body.Comment
	}

	if body.Active != nil {
		APIKey.Active = *body.Active
	}

	if body.Revoked != nil {
		APIKey.Revoked = *body.Revoked
	}

	if body.Ratelimit != nil {
		APIKey.Ratelimit = *body.Ratelimit
	}

	err = s.APIKeyStore.UpdateKey(APIKey)
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, APIKey)
}

func (s *APIServer) HandleDeleteKey(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	if !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedKey, err := strconv.Atoi(mux.Vars(r)["key_id"])
	if err != nil {
		return err
	}

	err = s.APIKeyStore.DeleteKey(requestedKey)
	if err != nil {
		return err
	}

	return SendEmptyResponse(w, http.StatusNoContent)
}

func GenerateJWT(key *common.APIKey) (string, error) {
	claims := &jwt.MapClaims{
		"id":       key.ID,
		"guild_id": key.GuildID,
	}

	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
