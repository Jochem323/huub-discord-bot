package api

import (
	"encoding/json"
	"huub-discord-bot/common"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) HandleGetKeywords(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusUnauthorized)
	}

	guildID := key.GuildID

	if key.Admin {
		guildID = r.URL.Query().Get("guild_id")
	}

	keywords, err := s.KeywordStore.GetKeywords(guildID)
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, keywords)
}

func (s *APIServer) HandleGetKeyword(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedKeyword, err := strconv.Atoi(mux.Vars(r)["keyword_id"])
	if err != nil {
		return err
	}

	keyword, err := s.KeywordStore.GetKeyword(requestedKeyword)
	if err != nil {
		return err
	}

	if keyword.GuildID != key.GuildID && !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	return SendDataResponse(w, http.StatusOK, keyword)
}

func (s *APIServer) HandleCreateKeyword(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	body := new(CreateKeywordRequest)
	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return err
	}

	guildId := key.GuildID

	if key.Admin {
		guildId = body.GuildID
	}

	keyword := common.Keyword{
		GuildID:  guildId,
		Keyword:  body.Keyword,
		Reaction: body.Reaction,
	}

	id, err := s.KeywordStore.AddKeyword(keyword)
	if err != nil {
		return err
	}

	keyword.ID = id

	return SendDataResponse(w, http.StatusCreated, keyword)

}

func (s *APIServer) HandleUpdateKeyword(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedKeyword, err := strconv.Atoi(mux.Vars(r)["keyword_id"])
	if err != nil {
		return err
	}

	body := new(UpdateKeywordRequest)
	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return err
	}

	keyword, err := s.KeywordStore.GetKeyword(requestedKeyword)
	if err != nil {
		return err
	}

	if keyword.GuildID != key.GuildID && !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	keyword.Reaction = body.Reaction

	err = s.KeywordStore.UpdateKeyword(keyword)
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, keyword)
}

func (s *APIServer) HandleDeleteKeyword(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedKeyword, err := strconv.Atoi(mux.Vars(r)["keyword_id"])
	if err != nil {
		return err
	}

	keyword, err := s.KeywordStore.GetKeyword(requestedKeyword)
	if err != nil {
		return err
	}

	if keyword.GuildID != key.GuildID && !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	err = s.KeywordStore.DeleteKeyword(requestedKeyword)
	if err != nil {
		return err
	}

	return SendEmptyResponse(w, http.StatusNoContent)
}
