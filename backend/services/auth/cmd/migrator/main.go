package main

import (
	"flag"
	"log"

	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/database"
)

func main() {
	fu := flag.Bool("up", false, "Apply all migrations")
	fd := flag.Int("down", 0, "Rollback N migrations")
	flag.Parse()

	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}

	m, err := database.NewMigrator(&cfg.Database, "./migrations")
	if err != nil {
		log.Fatalln("FATAL ->", err.Error())
	}
	defer m.Close()

	if *fu {
		if err := m.Up(); err != nil {
			log.Fatalln("FATAL ->", err.Error())
		}
	} else if *fd >= 0 {
		if err := m.Down(*fd); err != nil {
			log.Fatalln("FATAL ->", err.Error())
		}
	} else {
		log.Fatalln("FATAL -> failed to run migrations: no action specified")
	}
}
