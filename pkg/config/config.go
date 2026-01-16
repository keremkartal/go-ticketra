package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv     string `mapstructure:"APP_ENV"`
	IdentityPort string `mapstructure:"IDENTITY_SERVER_PORT"`
	EventPort string `mapstructure:"EVENT_SERVER_PORT"`
	BookingPort string `mapstructure:"BOOKING_SERVICE_PORT"`
	BookingServiceURL string `mapstructure:"BOOKING_SERVICE_URL"`
	// Postgres Settings
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`

	// Mongo Settings
	MongoUser     string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	MongoPassword string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	MongoDBName   string `mapstructure:"MONGO_DB"`
	MongoHost     string `mapstructure:"MONGO_HOST"`
	MongoPort     string `mapstructure:"MONGO_PORT"`

	// Redis Settings
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PASSWORD"`

	// RabbitMQ Settings
	RabbitUser string `mapstructure:"RABBITMQ_DEFAULT_USER"`
	RabbitPass string `mapstructure:"RABBITMQ_DEFAULT_PASS"`
	RabbitHost string `mapstructure:"RABBITMQ_HOST"`
	RabbitPort string `mapstructure:"RABBITMQ_PORT"`

	// JWT Settings
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiration int    `mapstructure:"JWT_EXPIRATION_HOURS"`

    // Gateway Settings
	GatewayPort       string `mapstructure:"GATEWAY_PORT"`
    IdentityServiceURL string `mapstructure:"IDENTITY_SERVICE_URL"`
    EventServiceURL    string `mapstructure:"EVENT_SERVICE_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}