#!/bin/bash

# Скрипт для создания админа
# Использование: ./create-admin.sh <email> <password>

if [ $# -ne 2 ]; then
    echo "Использование: $0 <email> <password>"
    echo "Пример: $0 admin@example.com MySecurePassword123"
    exit 1
fi

EMAIL=$1
PASSWORD=$2

# Проверяем наличие DATABASE_URL
if [ -z "$DATABASE_URL" ]; then
    echo "Ошибка: DATABASE_URL не установлен"
    echo "Установите: export DATABASE_URL='host=localhost dbname=restapi_prod user=postgres password=YOUR_PASSWORD sslmode=disable'"
    exit 1
fi

# Проверяем наличие Go
if ! command -v go &> /dev/null; then
    echo "Установка Go для создания админа..."
    # Можно использовать Docker для запуска Go скрипта
    echo "Используем альтернативный метод через SQL..."
    
    # Генерируем хеш пароля через Python (если доступен)
    if command -v python3 &> /dev/null; then
        HASH=$(python3 -c "
import bcrypt
import sys
password = sys.argv[1].encode('utf-8')
hash = bcrypt.hashpw(password, bcrypt.gensalt())
print(hash.decode('utf-8'))
" "$PASSWORD")
        
        # Вставляем админа в БД
        psql "$DATABASE_URL" -c "INSERT INTO admins (email, password_hash) VALUES ('$EMAIL', '$HASH') ON CONFLICT (email) DO NOTHING;"
        
        if [ $? -eq 0 ]; then
            echo "✅ Админ создан успешно!"
            echo "Email: $EMAIL"
        else
            echo "❌ Ошибка при создании админа"
        fi
    else
        echo "❌ Python3 не найден. Используйте метод через Go скрипт."
        exit 1
    fi
else
    # Используем Go скрипт
    go run -tags dev ./scripts/create-admin.go "$EMAIL" "$PASSWORD"
fi

