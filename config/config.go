package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Config *ApplicationConfig
var configFilePath string

type MySQLConfig struct {
	Hostname           string
	DBName             string
	Port               int
	Username           string
	Password           string
	Timezone           string
	ParseTime          bool
	MaxOpenConnections int
	MaxIdleConnections int
}

type ApplicationConfig struct {
	MySQL *MySQLConfig
}

func GetConfig() ApplicationConfig {
	if Config != nil {
		return *Config
	}
	err := InitConfig()
	if err != nil {
		panic(err)
	}
	return *Config
}

func InitConfig() error {
	if configFilePath == "" {
		return loadConfiguration()
	}

	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read env file: %w", err)
	}
	return loadConfiguration()
}

func loadConfiguration() error {
	viper.AutomaticEnv()

	config := &ApplicationConfig{
		MySQL: &MySQLConfig{
			Hostname:  viper.GetString("MYSQL_HOSTNAME"),
			DBName:    viper.GetString("MYSQL_DBNAME"),
			Port:      viper.GetInt("MYSQL_PORT"),
			Username:  viper.GetString("MYSQL_USERNAME"),
			Password:  viper.GetString("MYSQL_PASSWORD"),
			ParseTime: viper.GetBool("MYSQL_PARSE_TIME"),
		},
	}

	Config = config
	return nil
}

func Init(env string) {
	configFilePath = env
}
