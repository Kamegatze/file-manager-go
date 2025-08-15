package main

import (
	"file-manager/configuration"
	"log"
)

func main() {
	if err := configuration.Runner(configuration.NewDatasourceStarter, configuration.NewServerStarter); err != nil {
		log.Panic(err)
	}
}
