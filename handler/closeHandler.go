package handler

import (
	"os"
	"strconv"

	"github.com/apex/log"
)

func CloseHandler(totalExp int, totalClients int, log *log.Entry, c chan os.Signal) {
	defer close(c)
	tExp := strconv.Itoa(totalExp)
	tCli := strconv.Itoa(totalClients)
	log.Info("Total Expressions Calculated: " + tExp)
	log.Info("Total Clients Connected: " + tCli)
}
