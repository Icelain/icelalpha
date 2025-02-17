package cmd

import (
	"icealpha/internal/controllers"
	"icealpha/internal/database"
	"icealpha/internal/router"
	"icealpha/pkg/imglatex"
	"icealpha/pkg/inference"
	"log"
	"os"
)

func Execute() {

	loadEnvVars()
	flags := getFlags()

	db, err := database.CreatePostgresDriver(os.Getenv("POSTGRES_URL"))
	if err != nil {

		log.Fatal("Could not connect to postgres database: ", err)

	}

	srv := router.NewRouter()
	srvconfig := router.RouterConfig{

		Port:      flags.HttpPort,
		DB:        db,
		ImgLatex:  imglatex.NewImgLatex(os.Getenv("GROQ_API_KEY")),
		LLMClient: inference.NewClaudeLLMClient(os.Getenv("CLAUDE_API_KEY")),
	}

	srv.SetConfig(&srvconfig)
	controllers.HandleAll(srv)

	srv.Serve()

}
