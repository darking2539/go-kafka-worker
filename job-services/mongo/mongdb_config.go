package mongo

type MongoDBConfig struct {
	Username           string `env:"MONGODB_USERNAME" json:"username"`
	Password           string `env:"MONGODB_PASSWORD" json:"password"`
	Hostname           string `env:"MONGODB_HOSTNAME" json:"hostname"`
	Port               string `env:"MONGODB_PORT" json:"port"`
	Options            string `env:"MONGODB_OPTIONS" json:"options"`
	Timeout            int64  `env:"MONGODB_TIMEOUT" json:"timeout"`
	MaxPoolSize        int64  `env:"MONGODB_MAX_POOL_SIZE" json:"max_pool_size"`
	MinPoolSize        int64  `env:"MONGODB_MIN_POOL_SIZE" json:"min_pool_size"`
	MaxIdleConnections int64  `env:"MYSQL_MAX_IDLE_CONNECTIONS" json:"max_idle_connections"`
}
