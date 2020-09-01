package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config struct {
	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASS string
	DB_NAME string
	SECRET  string
	AWS_ACCESS_KEY_ID string
	AWS_SECRET_ACCESS_KEY string
	AWS_S3_BUCKET_NAME string
	FIREBASE_PRIVATEKEY string
	ITEMS_PER_PAGE int64
}

var config Config

func ReadConfig() {

	viper.SetConfigFile(".env")

	viper.AddConfigPath("../..")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

func GetConfig() Config {
	return config
}