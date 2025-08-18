package configuration

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gookit/config/v2"
)

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int32  `mapstructure:"port"`
}

func (server ServerConfig) Address() (string, error) {

	variableEmpty := []string{}

	if server.Host == "" {
		variableEmpty = append(variableEmpty, "Host: ''")
	}

	if server.Port == 0 {
		variableEmpty = append(variableEmpty, "Port: 0")
	}

	if len(variableEmpty) > 0 {
		return "", fmt.Errorf("several variable incorrect: %#v", variableEmpty)
	}

	return fmt.Sprintf("%s:%d", server.Host, server.Port), nil
}

func NewServer() (ServerConfig, error) {
	server := ServerConfig{}

	if err := InitConfig(); err != nil {
		return server, err
	}

	if err := config.BindStruct("server", &server); err != nil {
		return server, err
	}
	return server, nil
}

func NewServerStarter() (Starter, error) {
	return NewServer()
}

func (server ServerConfig) Run() error {
	router := gin.Default()

	server, err := NewServer()

	if err != nil {
		return err
	}

	address, err := server.Address()

	if err != nil {
		return err
	}

	router.Run(address)

	return nil
}

func (server ServerConfig) Close() error {
	return nil
}
