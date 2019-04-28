package handlers

import (
	"time"

	"github.com/dev-schueppchen/Kirby/internal/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-schueppchen/Kirby/internal/database"
)

type Message struct {
	db database.Middleware
}

func NewMessage(db database.Middleware) *Message {
	return &Message{
		db: db,
	}
}

func (h *Message) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Type != discordgo.MessageTypeDefault || e.Author.ID == s.State.User.ID {
		return
	}

	memb, err := s.GuildMember(e.GuildID, e.Author.ID)
	if err != nil {
		logger.Error("DISCORD :: failed getting guild member: %s", err.Error())
		return
	}

	msg := &database.Message{
		Attachment:       len(e.Attachments) > 0,
		Bot:              memb.User.Bot,
		ChannelID:        e.ChannelID,
		ContentLen:       len(e.Message.Content),
		GuildID:          e.GuildID,
		MentionedRoleIDs: e.MentionRoles,
		Mentions:         len(e.Mentions),
		RoleIDs:          memb.Roles,
		Timestamp:        time.Now(),
	}

	if err = h.db.AddMessage(msg); err != nil {
		logger.Error("DATABASE :: messageHandler: %s", err.Error())
	}
}
