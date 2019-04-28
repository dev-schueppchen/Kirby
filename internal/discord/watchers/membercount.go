package watchers

import (
	"time"

	"github.com/dev-schueppchen/Kirby/internal/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/dev-schueppchen/Kirby/internal/database"
)

const tickerDuration = 30 * time.Minute

// const tickerDuration = 5 * time.Second // for testing purposes

type MemberCount struct {
	ticker *time.Ticker
	s      *discordgo.Session
	db     database.Middleware
}

func NewMemberCount(s *discordgo.Session, db database.Middleware) *MemberCount {
	mc := &MemberCount{
		ticker: time.NewTicker(tickerDuration),
		s:      s,
		db:     db,
	}

	go func() {
		for {
			<-mc.ticker.C
			mc.onTick()
		}
	}()

	return mc
}

func (mc *MemberCount) onTick() {
	for _, guild := range mc.s.State.Guilds {
		go mc.pollRoleStatus(guild)
	}
}

func (mc *MemberCount) pollRoleStatus(guild *discordgo.Guild) {
	roleStatus := make(map[string]*database.MemberRoleStatus)
	for _, role := range guild.Roles {
		roleStatus[role.ID] = &database.MemberRoleStatus{
			RoleID: role.ID,
		}
	}

	ms := &database.MemberStatus{
		GuildID:   guild.ID,
		Timestamp: time.Now(),
	}

	for _, memb := range guild.Members {
		presence, err := mc.s.State.Presence(guild.ID, memb.User.ID)
		if err != nil {
			continue
		}

		for _, rID := range memb.Roles {
			rs := roleStatus[rID]
			switch presence.Status {
			case discordgo.StatusOnline:
				rs.Online++
			case discordgo.StatusIdle:
				rs.Away++
			case discordgo.StatusDoNotDisturb:
				rs.Dnd++
			default:
				rs.Offline++
			}
		}

		switch presence.Status {
		case discordgo.StatusOnline:
			ms.Online++
		case discordgo.StatusIdle:
			ms.Away++
		case discordgo.StatusDoNotDisturb:
			ms.Dnd++
		}
	}

	ms.Offline = guild.MemberCount - ms.Online - ms.Away - ms.Dnd

	mrsc := &database.MemberStatusRolesCollection{
		GuildID:   guild.ID,
		Timestamp: time.Now(),
		Roles:     make([]*database.MemberRoleStatus, len(roleStatus)),
	}

	var i int
	for _, v := range roleStatus {
		mrsc.Roles[i] = v
		i++
	}

	if err := mc.db.AddMembStatusRoles(mrsc); err != nil {
		logger.Error("DATABASE :: membRoleStatusCollection: %s", err.Error())
	}

	if err := mc.db.AddMembStatus(ms); err != nil {
		logger.Error("DATABASE :: membStatus: %s", err.Error())
	}
}
