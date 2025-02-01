package cmd

import (
	"icealpha/internal/controllers"
	"icealpha/internal/router"
)

func Execute() {

	flags := getFlags()

	srv := router.NewRouter()
	srvconfig := router.RouterConfig{

		Port: flags.HttpPort,
	}

	srv.SetConfig(&srvconfig)
	controllers.HandleAll(srv)

	srv.Serve()

}
