package config

type Config struct {
	HTTP HTTP
	DB   DB
}

type HTTP struct {
	Port int32
}

type DB struct{}

func Load() Config {
	return Config{}
}
