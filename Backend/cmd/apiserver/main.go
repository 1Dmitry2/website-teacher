package main

import (
	"backed-teacher/internal/app/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "configs/apiserver.toml", "Path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	
	// Переменные окружения имеют приоритет над конфигом
	if smtpHost := os.Getenv("SMTP_HOST"); smtpHost != "" {
		config.SMTP.Host = smtpHost
	}
	if smtpPort := os.Getenv("SMTP_PORT"); smtpPort != "" {
		if port, err := strconv.Atoi(smtpPort); err == nil {
			config.SMTP.Port = port
		}
	}
	if smtpUser := os.Getenv("SMTP_USERNAME"); smtpUser != "" {
		config.SMTP.Username = smtpUser
	}
	if smtpPass := os.Getenv("SMTP_PASSWORD"); smtpPass != "" {
		config.SMTP.Password = smtpPass
	}
	if smtpFrom := os.Getenv("SMTP_FROM"); smtpFrom != "" {
		config.SMTP.From = smtpFrom
	}
	
	// CORS origins from environment
	if corsOrigins := os.Getenv("CORS_ORIGINS"); corsOrigins != "" {
		origins := strings.Split(corsOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		config.CORSOrigins = origins
	}
	
	// Database URL from environment
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DataBaseURL = dbURL
	}
	
	// JWT Secret from environment
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	}
	
	// Admin reset URL from environment
	if adminResetURL := os.Getenv("ADMIN_RESET_URL"); adminResetURL != "" {
		config.Admin.ResetURL = adminResetURL
	}
	
	// User verification URL from environment
	if userVerificationURL := os.Getenv("USER_VERIFICATION_URL"); userVerificationURL != "" {
		config.User.VerificationURL = userVerificationURL
	}
	
	// User reset URL from environment
	if userResetURL := os.Getenv("USER_RESET_URL"); userResetURL != "" {
		config.User.ResetURL = userResetURL
	}
	
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
