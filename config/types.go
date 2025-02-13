package config

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"db"`
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
	Timezone string `json:"timezone"`
	SSLMode string `json:"ssl_mode"`
}

type ServerConfig struct {
	HttpPort          uint   `json:"http_port"`
	HttpHost          uint   `json:"http_host"`
	Secret            string `json:"secret"`
	AuthExpMinute     uint   `json:"auth_exp_min"`
	AuthRefreshMinute uint   `json:"auth_exp_refresh_min"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}
