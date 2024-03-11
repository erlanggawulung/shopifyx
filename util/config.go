package util

import "github.com/spf13/viper"

type Config struct {
	DBName            string `mapstructure:"DB_NAME"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBUsername        string `mapstructure:"DB_USERNAME"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	PromotheusAddress string `mapstructure:"PROMETHEUS_ADDRESS"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	BcryptSalt        string `mapstructure:"BCRYPT_SALT"`
	S3ID              string `mapstructure:"S3_ID"`
	S3SecretKey       string `mapstructure:"S3_SECRET_KEY"`
	S3BaseURL         string `mapstructure:"S3_BASE_URL"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
