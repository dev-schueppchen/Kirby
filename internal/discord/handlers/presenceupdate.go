package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/kirby/internal/database"
	"github.com/zekroTJA/kirby/internal/logger"
)

type PresenceUpdate struct {
	db         database.Middleware
	lastStatus map[string]discordgo.Status
}

func NewPresenceUpdate(db database.Middleware) *PresenceUpdate {
	return &PresenceUpdate{
		db:         db,
		lastStatus: make(map[string]discordgo.Status),
	}
}

func (h *PresenceUpdate) Handler(s *discordgo.Session, e *discordgo.PresenceUpdate) {
	lastStatus, ok := h.lastStatus[e.User.ID]
	if !ok {
		h.lastStatus[e.User.ID] = e.Status
		return
	}

	if lastStatus == e.Status {
		return
	}

	memb, err := s.GuildMember(e.GuildID, e.User.ID)
	if err != nil {
		logger.Error("DISCORD :: failed getting guild member: %s", err.Error())
		return
	}

	status := &database.StatusUpdate{
		Bot:       e.User.Bot,
		GuildID:   e.GuildID,
		NewStatus: e.Status,
		OldStatus: lastStatus,
		RoleIDs:   memb.Roles,
		Timestamp: time.Now(),
	}

	logger.Debug("%+v", status)
	h.lastStatus[e.User.ID] = e.Status

	if err = h.db.AddStatusUpdate(status); err != nil {
		logger.Error("DATABASE :: presenceUpdateHandler: %s", err.Error())
	}
}
