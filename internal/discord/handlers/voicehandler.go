package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-schueppchen/Kirby/internal/database"
	"github.com/dev-schueppchen/Kirby/internal/logger"
	"github.com/zekroTJA/timedmap"
)

const (
	voiceStateLifetime     = 3 * time.Hour
	voiceStateCleanupDelay = 15 * time.Minute
)

type Voice struct {
	db              database.Middleware
	lastVoiceStates *timedmap.TimedMap
}

func NewVoice(db database.Middleware) *Voice {
	return &Voice{
		db:              db,
		lastVoiceStates: timedmap.New(voiceStateCleanupDelay),
	}
}

func (h *Voice) Handler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	var oldVS *discordgo.VoiceState
	if h.lastVoiceStates.Contains(e.UserID) {
		oldVS, _ = h.lastVoiceStates.GetValue(e.UserID).(*discordgo.VoiceState)
	}

	newVS := e.VoiceState
	defer h.lastVoiceStates.Set(e.UserID, newVS, voiceStateLifetime)

	switch {

	// User left the voice channel
	case oldVS != nil && oldVS.ChannelID != "" && newVS.ChannelID == "":
		h.onVoiceChannelLeft(s, oldVS, newVS)

	// User moved to another channel
	case oldVS != nil && oldVS.ChannelID != "" && newVS.ChannelID != "" && oldVS.ChannelID != newVS.ChannelID:
		h.onVoiceChannelChange(s, oldVS, newVS)

	// User joins a channel
	case (oldVS == nil || oldVS.ChannelID == "") && newVS.ChannelID != "":
		h.onVoiceChannelJoined(s, oldVS, newVS)

	}
}

func (h *Voice) onVoiceChannelLeft(s *discordgo.Session, oldVS, newVS *discordgo.VoiceState) {
	h.addVoiceEvent(s, oldVS, database.VoiceLeft)
}

func (h *Voice) onVoiceChannelChange(s *discordgo.Session, oldVS, newVS *discordgo.VoiceState) {
	h.addVoiceEvent(s, newVS, database.VoiceMoved)
}

func (h *Voice) onVoiceChannelJoined(s *discordgo.Session, oldVS, newVS *discordgo.VoiceState) {
	h.addVoiceEvent(s, newVS, database.VoiceJoined)
}

func (h *Voice) addVoiceEvent(s *discordgo.Session, vc *discordgo.VoiceState, event database.VoiceEvent) {
	v := &database.Voice{
		GuildID:   vc.GuildID,
		Event:     event,
		ChannelID: vc.ChannelID,
		Timestamp: time.Now(),
	}

	memb, err := s.GuildMember(vc.GuildID, vc.UserID)
	if err != nil {
		logger.Error("DISCORD :: failed getting guild member: %s", err.Error())
		return
	}

	v.RoleIDs = memb.Roles

	if err := h.db.AddVoice(v); err != nil {
		logger.Error("DATABASE :: voiceHandler: %s", err.Error())
	}
}
