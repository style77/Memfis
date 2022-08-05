package main

import (
	api "Memfis/api"
	managers "Memfis/managers"
	"io"
	"log"
	"os"
)

func setupLog() {
	file, err := os.OpenFile("logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}

func main() {
	setupLog()

	config, err := managers.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	
	api.Run(config.ServerAddress)
}
