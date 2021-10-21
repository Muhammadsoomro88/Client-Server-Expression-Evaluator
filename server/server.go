package server

import (
	"net"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/cakemarketing/go-common/v5/settings"
)

func ConnectServer() net.Listener {
	settings.AddConfigPath("config")
	settings.SetConfigName("local")
	err := settings.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.SetHandler(cli.Default)
	log := log.WithFields(log.Fields{})

	dstream, err := net.Listen(settings.GetString("SERVER"), settings.GetString("TCP_PORT"))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Accept Conection on port 9000")
	return dstream
}
