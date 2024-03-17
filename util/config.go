package util

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBName            string `mapstructure:"DB_NAME"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBUsername        string `mapstructure:"DB_USERNAME"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	SSLMode           string `mapstructure:"SSL_MODE"`
	SSLRootCert       string `mapstructure:"SSL_ROOT_CERT"`
	PromotheusAddress string `mapstructure:"PROMETHEUS_ADDRESS"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	BcryptSalt        string `mapstructure:"BCRYPT_SALT"`
	S3ID              string `mapstructure:"S3_ID"`
	S3SecretKey       string `mapstructure:"S3_SECRET_KEY"`
	S3BucketName      string `mapstructure:"S3_BUCKET_NAME"`
}

func LoadConfig(path string) (Config, error) {
	config, err := LoadConfigFromFile(path)
	if os.Getenv("ENV") == "production" {
		return generateConfigFromEnvVars(), nil
	}
	return config, err
}

func LoadConfigFromFile(path string) (Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName("local")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}

func generateConfigFromEnvVars() Config {
	viper.AutomaticEnv()
	return Config{
		DBName:            viper.GetString("DB_NAME"),
		DBPort:            viper.GetString("DB_PORT"),
		DBHost:            viper.GetString("DB_HOST"),
		DBUsername:        viper.GetString("DB_USERNAME"),
		DBPassword:        viper.GetString("DB_PASSWORD"),
		SSLMode:           viper.GetString("SSL_MODE"),
		SSLRootCert:       viper.GetString("SSL_ROOT_CERT"),
		PromotheusAddress: viper.GetString("PROMOTHEUS_ADDRESS"),
		JWTSecret:         viper.GetString("JWT_SECRET"),
		BcryptSalt:        viper.GetString("BCRYPT_SALT"),
		S3ID:              viper.GetString("S3_ID"),
		S3SecretKey:       viper.GetString("S3_SECRET_KEY"),
		S3BucketName:      viper.GetString("S3_BUCKET_NAME"),
	}
}
