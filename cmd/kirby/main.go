package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/zekroTJA/kirby/internal/database"
	"github.com/zekroTJA/kirby/internal/discord"

	"github.com/zekroTJA/kirby/internal/config"
	"github.com/zekroTJA/kirby/internal/logger"
)

var (
	flagConfig   = flag.String("c", "config.yml", "config file location")
	flagLogLevel = flag.Int("loglvl", 4, "logger level [0 - critical, 5 - debug]")
)

func main() {
	flag.Parse()

	// LOGGER SETUP
	logger.Setup(`%{color}▶  %{level:.4s} %{id:03d}%{color:reset} %{message}`, *flagLogLevel)

	// CONMFIG PARSING
	cfg, err := config.Open(*flagConfig)
	if err != nil {
		logger.Fatal("CONFIG :: failed opening or creating config")
	}
	if cfg == nil {
		logger.Info("CONFIG :: config file was not found and was created at '%s'. "+
			"Enter your preferences and restart.", *flagConfig)
	}

	// DATABASE SETUP
	db := new(database.MongoDB)
	if err = db.Connect(cfg.MongoDB); err != nil {
		logger.Fatal("DATABASE :: failed connecting: %s", err.Error())
	}

	logger.Info("DATABASE :: connected")

	defer func() {
		logger.Info("DATABASE :: shutting down")
		db.Close()
	}()

	// DISCORD BOT SETUP
	dc, err := discord.New(cfg.Discord, db)
	if err != nil {
		logger.Fatal("DISCORD :: failed creating session: %s", err.Error())
	}

	go func() {
		if err := dc.OpenBlocking(); err != nil {
			logger.Fatal("DISCORD :: failed connecting to API: %s", err.Error())
		}
		logger.Info("DISCORD :: connected")
	}()

	defer func() {
		logger.Info("DISCORD :: shutting down")
		dc.Close()
	}()

	// Block main go routine untiul exit signal was detected
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
