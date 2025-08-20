package main

import (
	"log"

	"github.com/Kamegatze/file-manager-go/internal/configuration"
)

func main() {
	if err := configuration.Runner("config/config.yml", configuration.NewDatasourceStarter, configuration.NewServerStarter); err != nil {
		log.Panic(err)
	}
}
