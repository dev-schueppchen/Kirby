package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dev-schueppchen/Kirby/internal/database"
	"github.com/dev-schueppchen/Kirby/internal/discord/watchers"
	"github.com/dev-schueppchen/Kirby/internal/logger"
	"github.com/dev-schueppchen/Kirby/internal/static"
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
