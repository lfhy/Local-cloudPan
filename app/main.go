package main

import (
	"local-cloud-api/server"

	"github.com/lfhy/log"
)

func main() {
	log.SetLogLevel(log.DEBUGLevel)
	log.SetNoPrintCodeLine(false)
	server.Run()
}
