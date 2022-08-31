package config

var config = &GlobalConfig{
	Name:           "ServerApp",
	Host:           "0.0.0.0",
	Port:           2818,
	MaxConn:        1000,
	MaxPackageSize: 4096,
}

type GlobalConfig struct {
	Name           string
	Host           string
	Port           int
	MaxConn        int
	MaxPackageSize int
}

func GetGlobalConfig() *GlobalConfig {
	return config
}
