package config

type Config struct {
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
	Mail Mail   `mapstructure:"mail"`
	Msg  Msg    `mapstructure:"msg"`
}

type Mail struct {
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
