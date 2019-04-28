package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dev-schueppchen/Kirby/internal/config"
	"github.com/dev-schueppchen/Kirby/internal/database"
	"github.com/dev-schueppchen/Kirby/internal/discord/handlers"
)

type Discord struct {
	db  database.Middleware
	cfg *config.Discord

	s *discordgo.Session
}

func New(cfg *config.Discord, db database.Middleware) (d *Discord, err error) {
	d = new(Discord)

	d.db = db
	d.cfg = cfg

	if d.s, err = discordgo.New("Bot " + cfg.Token); err != nil {
		return
	}

	changeHandler := handlers.NewMemberChange(db)

	d.registerHandlers(
		handlers.NewReady(d.s, db).Handler,
		handlers.NewMessage(db).Handler,
		handlers.NewPresenceUpdate(db).Handler,
		handlers.NewReaction(db).Handler,
		handlers.NewVoice(db).Handler,
		changeHandler.HandlerAdd,
		changeHandler.HandlerRemove,
	)

	return
}

func (d *Discord) registerHandlers(handler ...interface{}) {
	for _, h := range handler {
		d.s.AddHandler(h)
	}
}

func (d *Discord) OpenBlocking() error {
	return d.s.Open()
}

func (d *Discord) Close() {
	d.s.Close()
}
