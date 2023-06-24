package main

import (
	"flag"
	"fmt"
	"log"
	"srbolab_cpc"
	"srbolab_cpc/db"
	"srbolab_cpc/server"
)

func main() {
	var conf srbolab_cpc.Config
	flag.StringVar(&conf.HTTP, "http", "127.0.0.1:4444", "HTTP address listen to")
	flag.StringVar(&conf.DbHost, "db-host", "127.0.0.1", "Database host")
	flag.IntVar(&conf.DbPort, "db_port", 5432, "Database port")
	flag.StringVar(&conf.DbName, "db-name", "srbolab_cpc", "Database name")
	flag.StringVar(&conf.DbUser, "db-user", "postgres", "Database user")
	flag.StringVar(&conf.DbPassword, "db-password", "p0stgres", "Database password")

	flag.Parse()

	if err := srbolab_cpc.LoadYamlConfig(&conf, "../config.yaml"); err != nil {
		log.Fatalf("Load configuration: %s", err)
	}

	err := db.Connect(conf)
	if err != nil {
		log.Fatalln(err)
	}

	server.RunServer(conf.HTTP)
	fmt.Println("cao")
}
