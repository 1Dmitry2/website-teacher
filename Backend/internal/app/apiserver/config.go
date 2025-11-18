package apiserver

type Config struct {
	BindAddr    string      `toml:"bind_addr"`
	LogLevel    string      `toml:"log_level"`
	DataBaseURL string      `toml:"database_url"`
	JWTSecret   string      `toml:"jwt_secret"`
	UploadDir   string      `toml:"upload_dir"`
	CORSOrigins []string    `toml:"cors_origins"`
	SMTP        SMTPConfig  `toml:"smtp"`
	Admin       AdminConfig `toml:"admin"`
	User        UserConfig  `toml:"user"`
}

type SMTPConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	From     string `toml:"from"`
}

type AdminConfig struct {
	ResetURL string `toml:"reset_url"`
}

type UserConfig struct {
	VerificationURL string `toml:"verification_url"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:  ":8080",
		LogLevel:  "debug",
		UploadDir: "./uploads",
		CORSOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		Admin: AdminConfig{
			ResetURL: "http://localhost:5173/admin/reset?token=",
		},
		User: UserConfig{
			VerificationURL: "http://localhost:5173/verify-email?token=",
		},
	}
}
