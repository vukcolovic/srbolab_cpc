package main

import (
	"flag"
	"log"
	"srbolab_cpc/config"
	"srbolab_cpc/db"
	"srbolab_cpc/handlers"
	"srbolab_cpc/server"
	"srbolab_cpc/service"
)

func main() {
	var conf config.Config
	flag.StringVar(&conf.HTTP, "http", "127.0.0.1:8001", "HTTP address listen to")
	flag.StringVar(&conf.DbHost, "db-host", "127.0.0.0", "Database host")
	flag.IntVar(&conf.DbPort, "db-port", 5431, "Database port")
	flag.StringVar(&conf.DbName, "db-name", "srbolab_cpc1", "Database name")
	flag.StringVar(&conf.DbUser, "db-user", "postgres1", "Database user")
	flag.StringVar(&conf.DbPassword, "db-password", "passw0rd1", "Database password")

	flag.Parse()

	if err := config.LoadYamlConfig(&conf, "config.yaml"); err != nil {
		log.Fatalf("Load configuration: %s", err)
	}

	err := db.Connect(conf)
	if err != nil {
		log.Fatalln(err)
	}

	handlers.CorporateIps = conf.CorporateIps
	service.Domain = conf.Domain
	service.RootPath = conf.RootPath

	server.RunServer(conf.HTTP)
}
