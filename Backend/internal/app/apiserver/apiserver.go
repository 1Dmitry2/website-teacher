package apiserver

import (
	"backed-teacher/internal/app/store/sqlstore"
	"backed-teacher/internal/mailer"
	"database/sql"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DataBaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)

	// Mailer опционален - если SMTP не настроен, приложение все равно работает
	// но письма не будут отправляться
	var mailerService *mailer.Mailer
	mailerService, err = mailer.New(mailer.Config{
		Host:     config.SMTP.Host,
		Port:     config.SMTP.Port,
		Username: config.SMTP.Username,
		Password: config.SMTP.Password,
		From:     config.SMTP.From,
	})
	if err != nil {
		// Логируем ошибку, но не останавливаем приложение
		// Это позволит работать без SMTP в режиме разработки
		mailerService = nil
	}

	srv := newServer(store, config.JWTSecret, mailerService, config.Admin.ResetURL, config.User.VerificationURL, config.UploadDir, config.CORSOrigins)
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
