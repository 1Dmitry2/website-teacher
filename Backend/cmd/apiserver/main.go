package main

import (
	"backed-teacher/internal/app/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"os"
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
		// Можно добавить парсинг порта, но для простоты оставим как есть
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
	
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
