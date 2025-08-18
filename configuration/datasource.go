package configuration

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gookit/config/v2"
	_ "github.com/lib/pq"
)

type DatasourceConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Host     string `mapstructure:"host"`
	Port     int32  `mapstructure:"port"`
	Driver   string `mapstructure:"driver"`
}

var datasource *sql.DB

func (datasourceConfig DatasourceConfig) CreateConnectionString() (string, error) {
	variableError := []string{}

	if datasourceConfig.Host == "" {
		log.Fatal("[DatasourceConfig] Host: ''")
		variableError = append(variableError, "Host: ''")
	}

	if datasourceConfig.Port == 0 {
		log.Fatal("[DatasourceConfig] Port: 0")
		variableError = append(variableError, "Port: 0")
	}

	if datasourceConfig.Username == "" {
		log.Fatal("[DatasourceConfig] Username: ''")
		variableError = append(variableError, "Username: ''")
	}

	if datasourceConfig.Password == "" {
		log.Fatal("[DatasourceConfig] Password: ''")
		variableError = append(variableError, "Password: ''")
	}

	if datasourceConfig.Database == "" {
		log.Fatal("[DatasourceConfig] Database: ''")
		variableError = append(variableError, "Database: ''")
	}

	if len(variableError) > 0 {
		return "", fmt.Errorf("several variable incorrect: %#v", variableError)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.Host, datasourceConfig.Port, datasourceConfig.Database), nil
}

func NewDatasource() (DatasourceConfig, error) {
	datasource := DatasourceConfig{}

	if err := InitConfig(); err != nil {
		log.Fatal(err)
		return datasource, err
	}

	if err := config.BindStruct("datasource", &datasource); err != nil {
		log.Fatal(err)
		return datasource, err
	}
	return datasource, nil
}

func NewDatasourceStarter() (Starter, error) {
	return NewDatasource()
}

func (datasourceConfig DatasourceConfig) OpenDatabase() (*sql.DB, error) {

	connection, err := datasourceConfig.CreateConnectionString()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return sql.Open(datasourceConfig.Driver, connection)
}

func (datasourceConfig DatasourceConfig) Migration(migrationPath string) error {
	connectionString, err := datasourceConfig.CreateConnectionString()

	if err != nil {
		log.Fatal(err)
		return err
	}
	m, err := migrate.New(migrationPath, connectionString)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
		return err
	}

	return nil
}

func (datasourceConfig DatasourceConfig) Run() error {
	if err := datasourceConfig.Migration("file:./migrations"); err != nil {
		log.Fatal(err)
		return err
	}
	return RunDatabase(datasourceConfig)
}

func (datasourceConfig DatasourceConfig) Close() error {
	if err := datasource.Close(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func RunDatabase(datasourceConfig DatasourceConfig) error {

	db, err := datasourceConfig.OpenDatabase()

	if err != nil {
		log.Fatal(err)
		return err
	}

	datasource = db

	return nil
}

func GetDB() (*sql.DB, error) {
	if datasource == nil {
		log.Fatal("sql.DB is nil")
		return nil, errors.New("sql.DB is nil")
	}
	return datasource, nil
}
