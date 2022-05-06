package config

type ServerConfig struct {
	Name        string      `mapstructure:"name"`
	Host        string      `mapstructure:"host"`
	Port        int         `mapstructure:"port"`
	LogsAddress string      `mapstructure:"logsAddress"`
	Mysqlinfo   MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}
