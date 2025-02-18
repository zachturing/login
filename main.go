package main

import (
	"flag"
	"fmt"

	"github.com/newdee/aipaper-util/log"
)

var (
	port = flag.Int("port", 30001, "port")
)

func main() {
	flag.Parse()

	if err := initService(); err != nil {
		log.Errorf("initService error: %v", err)
		return
	}

	router := initRoute()

	addr := fmt.Sprintf("0.0.0.0:%v", *port)
	log.Debugf("start server at ", addr)
	router.Run(addr)
}
