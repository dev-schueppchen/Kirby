package database

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type MembChangeEvent int

type VoiceEvent int

const (
	ChangeJoined MembChangeEvent = iota
	ChangeLeft
)

const (
	VoiceJoined VoiceEvent = iota
	VoiceMoved
	VoiceLeft
)

type Message struct {
	GuildID          string
	ChannelID        string
	RoleIDs          []string
	Bot              bool
	ContentLen       int
	Timestamp        time.Time
	MentionedRoleIDs []string
	Mentions         int
	Attachment       bool
}

type StatusUpdate struct {
	GuildID   string
	RoleIDs   []string
	Bot       bool
	OldStatus discordgo.Status
	NewStatus discordgo.Status
	Timestamp time.Time
}

type MemberChange struct {
	GuildID   string
	Timestamp time.Time
	RoleIDs   []string
	Bot       bool
	Event     MembChangeEvent
}

type MemberRoleStatus struct {
	RoleID  string
	Online  int
	Offline int
	Dnd     int
	Away    int
}

type MemberStatusRolesCollection struct {
	GuildID   string
	Timestamp time.Time
	Roles     []*MemberRoleStatus
}

type MemberStatus struct {
	GuildID   string
	Timestamp time.Time
	Online    int
	Offline   int
	Dnd       int
	Away      int
}

type Reaction struct {
	GuildID    string
	Timestamp  time.Time
	ChannelID  string
	RoleIDs    []string
	Emoji      string
	EmojiID    string
	ContentLen int
	Bot        bool
}

type Voice struct {
	GuildID   string
	ChannelID string
	Timestamp time.Time
	RoleIDs   []string
	Event     VoiceEvent
}

type Middleware interface {
	Connect(params interface{}) error
	Close()

	AddMessage(msg *Message) error
	AddStatusUpdate(su *StatusUpdate) error
	AddMembChange(mc *MemberChange) error
	AddMembStatus(ms *MemberStatus) error
	AddMembStatusRoles(mrsc *MemberStatusRolesCollection) error
	AddReaction(r *Reaction) error
	AddVoice(v *Voice) error
}
