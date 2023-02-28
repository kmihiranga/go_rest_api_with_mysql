package config

type Config struct {
	AppConfig *AppConfig
	DBConfig  *DBConfig
}

type AppConfig struct {
	Name                string `yaml:"name"`
	Host                string `yaml:"host"`
	Port                int    `yaml:"port"`
	ReadTimeout         uint   `yaml:"read-timeout"`
	WriteTimeout        uint   `yaml:"write-timeout"`
	ShutdownWaitTimeout uint   `yaml:"shutdown-wait-timeout"`
	EncryptKey          string `yaml:"enc-key"`
	BaseURL             string `yaml:"base-url"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
