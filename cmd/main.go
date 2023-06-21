package main

import (
	"flag"
	"log"
	"srbolab_cpc/server"
)

func main() {
	var conf Config
	flag.StringVar(&conf.HTTP, "http", "http:127.0.0.1", "HTTP address listen to")
	flag.IntVar(&conf.HTTPPort, "http_port", 4444, "HTTP port")
	flag.StringVar(&conf.DbHost, "db-host", "", "Database host")
	flag.StringVar(&conf.DbName, "db-name", "", "Database name")
	flag.StringVar(&conf.DbUser, "db-user", "", "Database user")
	flag.StringVar(&conf.DbPassword, "db-password", "", "Database password")
	flag.Parse()

	if err := loadYamlConfig(&conf, "config.yaml"); err != nil {
		log.Fatalf("Load configuration: %s", err)
	}

	server.RunServer(conf.HTTP, 10)
}
