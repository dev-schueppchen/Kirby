package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/kirby/internal/database"
	"github.com/zekroTJA/kirby/internal/logger"
)

type MemberChange struct {
	db database.Middleware
}

func NewMemberChange(db database.Middleware) *MemberChange {
	return &MemberChange{
		db: db,
	}
}

func (h *MemberChange) HandlerAdd(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	c := &database.MemberChange{
		Bot:       e.Member.User.Bot,
		Event:     database.ChangeJoined,
		GuildID:   e.GuildID,
		RoleIDs:   e.Member.Roles,
		Timestamp: time.Now(),
	}

	h.addDBEntry(c)
}

func (h *MemberChange) HandlerRemove(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
	c := &database.MemberChange{
		Bot:       e.Member.User.Bot,
		Event:     database.ChangeLeft,
		GuildID:   e.GuildID,
		RoleIDs:   e.Member.Roles,
		Timestamp: time.Now(),
	}

	h.addDBEntry(c)
}

func (h *MemberChange) addDBEntry(c *database.MemberChange) {
	if err := h.db.AddMembChange(c); err != nil {
		logger.Error("DATABASE :: memberChangeHandler: %s", err.Error())
	}
}
