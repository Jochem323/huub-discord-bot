package api

import (
	"encoding/json"
	"huub-discord-bot/common"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) HandleGetGuilds(w http.ResponseWriter, r *http.Request) error {
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

	guilds, err := s.GuildStore.GetGuilds()
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, guilds)
}

func (s *APIServer) HandleGetGuild(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedGuild := mux.Vars(r)["guild_id"]

	if requestedGuild != key.GuildID && !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	guild, err := s.GuildStore.GetGuild(requestedGuild)
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, guild)
}

func (s *APIServer) HandleCreateGuild(w http.ResponseWriter, r *http.Request) error {
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

	body := new(CreateGuildRequest)
	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return err
	}

	guild := common.Guild{
		ID:     body.ID,
		Prefix: body.Prefix,
	}

	err = s.GuildStore.AddGuild(guild)
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusCreated, guild)

}

func (s *APIServer) HandleUpdateGuild(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedGuild := mux.Vars(r)["guild_id"]

	if requestedGuild != key.GuildID && !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	body := new(UpdateGuildRequest)
	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return err
	}

	guild, err := s.GuildStore.GetGuild(requestedGuild)
	if err != nil {
		return err
	}

	if body.Prefix != "" {
		guild.Prefix = body.Prefix
	}

	err = s.GuildStore.UpdateGuild(guild)
	if err != nil {
		return err
	}

	return SendDataResponse(w, http.StatusOK, guild)
}

func (s *APIServer) HandleDeleteGuild(w http.ResponseWriter, r *http.Request) error {
	key, err := s.GetKeyFromRequest(r)
	if err != nil {
		return err
	}

	if !s.GetKeyValidity(key) {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	requestedGuild := mux.Vars(r)["guild_id"]

	if requestedGuild != key.GuildID && !key.Admin {
		return SendEmptyResponse(w, http.StatusForbidden)
	}

	err = s.GuildStore.DeleteGuild(requestedGuild)
	if err != nil {
		return err
	}

	return SendEmptyResponse(w, http.StatusNoContent)
}
