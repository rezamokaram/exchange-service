package config

type AConfig struct {
	DB     DBConfig     `json:"postgres"`
	Server ServerConfig `json:"app"`
	Redis  RedisConfig  `json:"redis"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"db"`
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ServerConfig struct {
	HttpPort          uint   `json:"port"`
	Secret            string `json:"secret"`
	AuthExpMinute     uint   `json:"auth_exp_min"`
	AuthRefreshMinute uint   `json:"auth_exp_refresh_min"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}