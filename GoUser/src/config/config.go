package config

type Config struct {
	Name     string   `mapstructure:"name"`
	Msg      Msg      `mapstructure:"msg"`
	Postgres Postgres `mapstructure:"postgres"`
	Register Redis    `mapstructure:"register"`
	Validate Redis    `mapstructure:"validate"`
}

type Msg struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Index    int    `mapstructure:"index"`
}
