package settings

type Config struct {
	MySql  MySqlSetting  `mapstructure:"mysql"`
	Server ServerSetting `mapstructure:"server"`
	Logger LoggerSetting `mapstructure:"logger"`
	Redis  RedisSetting  `mapstructure:"redis"`
	Smtp   SMTPSetting   `mapstructure:"smtp"`
	Kafka  KafkaSetting  `mapstructure:"kafka"`
	Jwt    JwtSetting    `mapstructure:"jwt"`
}

type MySqlSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"db_name"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxIdle         int    `mapstructure:"max_idle_conns"`
	MaxOpen         int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type ServerSetting struct {
	Port    int      `mapstructure:"port"`
	Cors    []string `mapstructure:"cors"`
	Version string   `mapstructure:"version"`
}

type LoggerSetting struct {
	Level       string `mapstructure:"level"`
	FileLogPath string `mapstructure:"file_log_path"`
	MaxSize     int    `mapstructure:"max_size"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      int    `mapstructure:"max_age"`
	Compress    bool   `mapstructure:"compress"`
}

type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type SMTPSetting struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	Secure bool   `mapstructure:"secure"`
}

type KafkaSetting struct {
	Host string `mapstructure:"host"`
}

type JwtSetting struct {
	AccessKey  string `mapstructure:"access_key"`
	RefreshKey string `mapstructure:"refresh_key"`
	AccessExp  string `mapstructure:"access_exp"`
	RefreshExp string `mapstructure:"refresh_exp"`
}
