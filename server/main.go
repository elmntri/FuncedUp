package main

import (
	"funcedup/internal/schema"
	"funcedup/pkg/config"
	"funcedup/pkg/logger"
	"funcedup/pkg/pgconn"
	"funcedup/pkg/server"

	"go.uber.org/fx"
)

func init() {
	//! CONFIG PRECEDENCE: OVERRIDE > ENV > CONFIG FILE > FALLBACK

	// *OPTIONAL* OVERRIDE GLOBAL LOG LEVEL, INTENDED FOR DEVELOPMENT
	// config.OverrideGlobalLogLevel("debug")

	config.SetUpConfig("SERVER", "yaml", "./")
}

func main() {
	app := fx.New(
		//* Modules ---------------------------------------------------------------
		logger.InjectModule("logger"),
		pgconn.InjectModule("database"),
		server.InjectModule("server"),
		//* Domains ---------------------------------------------------------------

		//* Migration -------------------------------------------------------------
		fx.Invoke(func(m *pgconn.Module) {
			m.ApplySchema(
				true,
				schema.User{},
				schema.Content{},
				schema.Discussion{},
				schema.DiscusionReply{},
				schema.Note{},
				schema.NoteReply{},
				schema.Tag{},
				schema.ContentTag{},
			)
		}),
		//* fx logs ---------------------------------------------------------------
		fx.NopLogger,
	)
	app.Run()
}
