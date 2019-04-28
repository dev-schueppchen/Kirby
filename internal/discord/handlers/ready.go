package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/kirby/internal/database"
	"github.com/zekroTJA/kirby/internal/discord/watchers"
	"github.com/zekroTJA/kirby/internal/logger"
	"github.com/zekroTJA/kirby/internal/static"
)

type Ready struct {
	s  *discordgo.Session
	db database.Middleware
}

func NewReady(s *discordgo.Session, db database.Middleware) *Ready {
	return &Ready{
		s:  s,
		db: db,
	}
}

func (h *Ready) Handler(s *discordgo.Session, e *discordgo.Ready) {
	logger.Info("DISCORD :: ready as user %s (%s)", s.State.User.String(), s.State.User.ID)
	logger.Info("DISCORD :: invite: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		e.User.ID, static.InvitePermission)

	h.startWatchers()
}

func (h *Ready) startWatchers() {
	watchers.NewMemberCount(h.s, h.db)
}
