package postgres

type Config struct {
	URL      string `yaml:"url" env:"PG_URL" env-required:"true"`
	MaxConns int32  `yaml:"max_conns" env:"PG_MAX_CONNS" env-default:"1"`
}
