package postgres

type Config struct {
	Database string `config:"DB_NAME" yaml:"database" validate:"required"`
	User     string `config:"DB_USER" yaml:"user" validate:"required"`
	Password string `config:"DB_PASSWORD" yaml:"password" validate:"required"`
	Host     string `config:"DB_HOST" yaml:"host" validate:"required"`
	Port     int    `config:"DB_PORT" yaml:"port" validate:"required"`
	Retries  int    `config:"DB_CONNECT_RETRY" yaml:"retries"`
	PoolSize int    `config:"DB_POOL_SIZE" yaml:"pool_size" validate:"required"`
}
