package main

import (
	"fmt"
	"local-cloud-api/conf"
	"local-cloud-api/server"
	"net"

	"github.com/lfhy/log"
)

func main() {
	if conf.Port == "auto" {
		okPort := 9527
		for {
			l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", okPort))
			if err != nil {
				okPort++
			} else {
				l.Close()
				break
			}
		}
		conf.Port = fmt.Sprint(okPort)
	}
	log.SetLogLevel(log.DEBUGLevel)
	log.SetNoPrintCodeLine(false)
	server.Run()
	if conf.Gui {
		server.RunGui()
	} else {
		select {}
	}
}
