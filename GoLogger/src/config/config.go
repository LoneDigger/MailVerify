package config

type Config struct {
	Name  string `mapstructure:"name"`
	Msg   Msg    `mapstructure:"msg"`
	Mongo Mongo  `mapstructure:"mongo"`
}

type Mongo struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Msg struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
