package handlers

import (
	"time"

	"github.com/dev-schueppchen/Kirby/internal/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-schueppchen/Kirby/internal/database"
)

type Reaction struct {
	db database.Middleware
}

func NewReaction(db database.Middleware) *Reaction {
	return &Reaction{
		db: db,
	}
}

func (h *Reaction) Handler(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	r := &database.Reaction{
		GuildID:   e.GuildID,
		Timestamp: time.Now(),
		ChannelID: e.ChannelID,
		Emoji:     e.Emoji.Name,
		EmojiID:   e.Emoji.ID,
	}

	memb, err := s.GuildMember(e.GuildID, e.UserID)
	if err != nil {
		logger.Error("DISCORD :: failed getting guild member: %s", err.Error())
		return
	}

	msg, err := s.ChannelMessage(e.ChannelID, e.MessageID)
	if err != nil {
		logger.Error("DISCORD :: failed getting message: %s", err.Error())
		return
	}

	r.Bot = memb.User.Bot
	r.RoleIDs = memb.Roles
	r.ContentLen = len(msg.Content)

	if err = h.db.AddReaction(r); err != nil {
		logger.Error("DATABASE :: reactionHandler: %s", err.Error())
	}
}
