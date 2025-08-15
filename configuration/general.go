package configuration

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var IsInit = false

func InitConfig() error {
	return InitConfigWithArgs(yaml.Driver, "resources/config.yml")
}

func InitConfigWithArgs(driver *config.StdDriver, path string) error {
	if !IsInit {
		config.WithOptions(config.ParseEnv)

		config.AddDriver(driver)

		if err := config.LoadFiles(path); err != nil {
			return err
		}
		IsInit = true
	}
	return nil
}
