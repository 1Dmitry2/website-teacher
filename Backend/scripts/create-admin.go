package main

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store/sqlstore"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	email := flag.String("email", "", "Admin email")
	password := flag.String("password", "", "Admin password")
	flag.Parse()

	if *email == "" || *password == "" {
		if len(os.Args) >= 3 {
			*email = os.Args[1]
			*password = os.Args[2]
		} else {
			fmt.Println("Использование: go run create-admin.go -email <email> -password <password>")
			fmt.Println("Или: go run create-admin.go <email> <password>")
			os.Exit(1)
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не установлен")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	store := sqlstore.New(db)

	admin := &model.Admin{
		Email:    *email,
		Password: *password,
	}

	if err := store.Admin().Create(admin); err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		if err.Error() == "pq: duplicate key value violates unique constraint \"admins_email_key\"" {
			fmt.Println("Админ с таким email уже существует")
		}
		os.Exit(1)
	}

	fmt.Println("✅ Админ создан успешно!")
	fmt.Printf("Email: %s\n", admin.Email)
	fmt.Printf("ID: %d\n", admin.ID)
}

