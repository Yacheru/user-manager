package config

import (
	"github.com/spf13/viper"
	"user-manager/init/logger"
	"user-manager/pkg/constants"
)

var CFG Config

type Config struct {
	ApiDebug bool   `mapstructure:"API_DEBUG"`
	ApiPort  int    `mapstructure:"API_PORT"`
	ApiEntry string `mapstructure:"API_ENTRY"`

	PostgresqlDSN string `mapstructure:"POSTGRESQL_DSN"`

	MongodbURI        string `mapstructure:"MONGODB_URI"`
	MongodbDatabase   string `mapstructure:"MONGODB_DATABASE"`
	MongodbCollection string `mapstructure:"MONGODB_COLLECTION"`
}

func InitConfig() error {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(err.Error(), constants.ConfigCategory)
		return err
	}

	if err := viper.Unmarshal(&CFG); err != nil {
		logger.Error(err.Error(), constants.ConfigCategory)
		return err
	}

	return nil
}
