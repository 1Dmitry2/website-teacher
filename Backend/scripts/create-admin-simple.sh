#!/bin/bash

# Простой скрипт для создания админа через SQL
# Использование: ./create-admin-simple.sh <email> <password>

if [ $# -ne 2 ]; then
    echo "Использование: $0 <email> <password>"
    exit 1
fi

EMAIL=$1
PASSWORD=$2

# Используем Python для генерации bcrypt хеша
if command -v python3 &> /dev/null; then
    HASH=$(python3 << EOF
import bcrypt
password = b'$PASSWORD'
hash = bcrypt.hashpw(password, bcrypt.gensalt())
print(hash.decode('utf-8'))
EOF
)
    
    # Вставляем в БД через docker-compose
    docker-compose exec -T postgres psql -U postgres -d restapi_prod << SQL
INSERT INTO admins (email, password_hash) 
VALUES ('$EMAIL', '$HASH') 
ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash;
SELECT 'Админ создан: ' || email FROM admins WHERE email = '$EMAIL';
SQL

elif command -v docker &> /dev/null; then
    # Используем временный контейнер с Python
    HASH=$(docker run --rm python:3-alpine python3 -c "
import bcrypt
password = b'$PASSWORD'
hash = bcrypt.hashpw(password, bcrypt.gensalt())
print(hash.decode('utf-8'))
")
    
    docker-compose exec -T postgres psql -U postgres -d restapi_prod << SQL
INSERT INTO admins (email, password_hash) 
VALUES ('$EMAIL', '$HASH') 
ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash;
SELECT 'Админ создан: ' || email FROM admins WHERE email = '$EMAIL';
SQL
else
    echo "Ошибка: Python3 или Docker не найдены"
    exit 1
fi

echo "✅ Админ создан!"
echo "Email: $EMAIL"

