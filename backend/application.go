package main

import (
	"flag"
	"fmt"
	"os"

	config "github.com/sriharivishnu/shopify-challenge/config"
	db "github.com/sriharivishnu/shopify-challenge/external"
	"github.com/sriharivishnu/shopify-challenge/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	config.PopulateConfig()

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	db.Init()
	server.Init()
}
