package main

import (
	"file-manager/internal/configuration"
	"log"
)

func main() {
	if err := configuration.Runner("config/config.yml", configuration.NewDatasourceStarter, configuration.NewServerStarter); err != nil {
		log.Panic(err)
	}
}
