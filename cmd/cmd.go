package cmd

import (
	"icealpha/internal/controllers"
	"icealpha/internal/database"
	"icealpha/internal/router"
	"log"
	"os"
)

func Execute() {

	loadEnvVars()
	flags := getFlags()

	db, err := database.CreatePostgresDriver(os.Getenv("POSTGRES_URL"))
	if err != nil {

		log.Fatal("Could not connect to postgres db")

	}

	srv := router.NewRouter()
	srvconfig := router.RouterConfig{

		Port: flags.HttpPort,
		DB:   db,
	}

	srv.SetConfig(&srvconfig)
	controllers.HandleAll(srv)

	srv.Serve()

}
