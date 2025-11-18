#!/bin/bash

# Создание админа через временный Go контейнер
# Использование: ./create-admin-docker.sh <email> <password>

if [ $# -ne 2 ]; then
    echo "Использование: $0 <email> <password>"
    exit 1
fi

EMAIL=$1
PASSWORD=$2

echo "Создание админа через Go контейнер..."

# Используем временный Go контейнер для генерации правильного хеша
cd /opt/website-teacher/Backend

# Копируем необходимые файлы во временную директорию
TMP_DIR=$(mktemp -d)
cp -r internal $TMP_DIR/
cp go.mod go.sum $TMP_DIR/
mkdir -p $TMP_DIR/scripts
cat > $TMP_DIR/scripts/gen-hash.go << 'GOHASH'
package main

import (
	"fmt"
	"os"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <password>\n", os.Args[0])
		os.Exit(1)
	}
	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(string(hash))
}
GOHASH

# Генерируем хеш через Go контейнер
HASH=$(docker run --rm -v "$TMP_DIR:/app" -w /app golang:1.24-alpine sh -c "
cd /app
go mod download
go run scripts/gen-hash.go '$PASSWORD'
")

# Очищаем временную директорию
rm -rf $TMP_DIR

if [ -z "$HASH" ]; then
    echo "❌ Ошибка при генерации хеша"
    exit 1
fi

# Вставляем админа в БД
docker-compose exec -T postgres psql -U postgres -d restapi_prod << SQL
INSERT INTO admins (email, password_hash) 
VALUES ('$EMAIL', '$HASH') 
ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash;
SELECT '✅ Админ создан: ' || email FROM admins WHERE email = '$EMAIL';
SQL

echo ""
echo "✅ Готово! Админ создан:"
echo "Email: $EMAIL"

