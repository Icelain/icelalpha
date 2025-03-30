package cmd

import (
	"icealpha/internal/controllers"
	"icealpha/internal/controllers/jwtauth"
	"icealpha/internal/database"
	"icealpha/internal/router"
	"icealpha/pkg/imglatex"
	"icealpha/pkg/inference"
	"log"
	"os"
	"sync"
)

func Execute() {

	loadEnvVars()
	flags := getFlags()

	db, err := database.CreatePostgresDriver(os.Getenv("POSTGRES_URL"))
	if err != nil {

		log.Fatal("Could not connect to postgres database: ", err)

	}

	log.Println("Established a connection with the database")

	srv := router.NewRouter()
	srvconfig := router.RouterConfig{

		Port:        flags.HttpPort,
		DB:          db,
		ImgLatex:    imglatex.NewImgLatex(os.Getenv("GROQ_API_KEY")),
		LLMClient:   inference.NewClaudeLLMClient(os.Getenv("CLAUDE_API_KEY")),
		JWTSession:  jwtauth.NewJWTSession(os.Getenv("SESSION_KEY")),
		CreditCache: &sync.Map{},
	}

	srv.SetConfig(&srvconfig)
	controllers.HandleAll(srv)

	srv.Serve()

}
