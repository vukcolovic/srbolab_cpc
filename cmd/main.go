package main

import (
	"flag"
	"log"
	"srbolab_cpc"
	"srbolab_cpc/db"
	"srbolab_cpc/handlers"
	"srbolab_cpc/server"
)

func main() {
	var conf srbolab_cpc.Config
	flag.StringVar(&conf.HTTP, "http", "127.0.0.1:8000", "HTTP address listen to")
	flag.StringVar(&conf.DbHost, "db-host", "127.0.0.1", "Database host")
	flag.IntVar(&conf.DbPort, "db_port", 5432, "Database port")
	flag.StringVar(&conf.DbName, "db-name", "srbolab_cpc", "Database name")
	flag.StringVar(&conf.DbUser, "db-user", "postgres", "Database user")
	flag.StringVar(&conf.DbPassword, "db-password", "password", "Database password")

	flag.Parse()

	if err := srbolab_cpc.LoadYamlConfig(&conf, "config.yaml"); err != nil {
		log.Fatalf("Load configuration: %s", err)
	}

	err := db.Connect(conf)
	if err != nil {
		log.Fatalln(err)
	}

	handlers.CorporateIps = conf.CorporateIps

	server.RunServer(conf.HTTP)
}
