package config

type Properties struct {
	Port           string `env:"MY_APP_PORT" env-default:"8080"`
	Host           string `env:"HOST" env-default:"localhost"`
	DBHost         string `env:"DB_HOST" env-default:"localhost"`
	DBPort         string `env:"DB_PORT" env-default:"27017"`
	DBName         string `env:"DB_NAME" env-default:"tronics"`
	CollectionName string `env:"COLLECTION_NAME" env-default:"products"`
}
