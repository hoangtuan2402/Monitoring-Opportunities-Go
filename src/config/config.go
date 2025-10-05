package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBName     string `mapstructure:"DB_NAME"`
	MongoURI   string `mapstructure:"MONGO_URI"`
	EthRPCURL  string `mapstructure:"ETH_RPC_URL"`
}

var envs = []string{
	"DB_NAME", "MONGO_URI", "ETH_RPC_URL",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
