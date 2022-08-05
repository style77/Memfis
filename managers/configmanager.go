package managers

import (
	"github.com/spf13/viper"
)

func LoadConfig() (config config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func WriteToConfig(option, value string) {
	viper.Set(option, value)
	viper.WriteConfig()
}

func CheckForSecret() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	if config.SaveSecret {
		WriteToConfig("SECRET_KEY", SecretKey)
	}
}
