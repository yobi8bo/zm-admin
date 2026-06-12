package config

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Log       LogConfig       `mapstructure:"log"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Storage   StorageConfig   `mapstructure:"storage"`
}

type ServerConfig struct {
	Port               int    `mapstructure:"port"`
	Mode               string `mapstructure:"mode"`
	JwtSecret          string `mapstructure:"jwt_secret"`
	AccessTokenExpire  int64  `mapstructure:"access_token_expire"`
	RefreshTokenExpire int64  `mapstructure:"refresh_token_expire"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LogLevel     string `mapstructure:"log_level"`
}

type RedisConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Password      string `mapstructure:"password"`
	DB            int    `mapstructure:"db"`
	PoolSize      int    `mapstructure:"pool_size"`
	KeyPrefix     string `mapstructure:"key_prefix"`
	CacheTTL      int64  `mapstructure:"cache_ttl"`
	CaptchaExpire int64  `mapstructure:"captcha_expire"`
	DialTimeout   int64  `mapstructure:"dial_timeout"`
	ReadTimeout   int64  `mapstructure:"read_timeout"`
	WriteTimeout  int64  `mapstructure:"write_timeout"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type RateLimitConfig struct {
	Rate           float64 `mapstructure:"rate"`
	Burst          int     `mapstructure:"burst"`
	LoginRate      float64 `mapstructure:"login_rate"`
	LoginBurst     int     `mapstructure:"login_burst"`
	CaptchaRate    float64 `mapstructure:"captcha_rate"`
	CaptchaBurst   int     `mapstructure:"captcha_burst"`
	SensitiveRate  float64 `mapstructure:"sensitive_rate"`
	SensitiveBurst int     `mapstructure:"sensitive_burst"`
}

type StorageConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}
