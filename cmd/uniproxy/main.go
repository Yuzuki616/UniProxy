package main

import (
	"UniProxy/router"
	"flag"
	"log"
)

var host = flag.String("host", "127.0.0.1", "host")
var port = flag.Int("port", 11451, "port")

func main() {
	flag.Parse()
	router.Init()
	if err := router.Start(*host, *port); err != nil {
		log.Fatalln("start error:", err)
	}
}
