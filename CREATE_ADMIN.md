# Создание админа

В проекте нет автоматического создания админа. Нужно создать его вручную.

## Способ 1: Через скрипт (рекомендуется)

На сервере выполните:

```bash
cd /opt/website-teacher/Backend

# Установите DATABASE_URL
export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=ql1OdITTGh2023woLJToH3MaUGWYEMFH sslmode=disable"

# Создайте админа
go run scripts/create-admin.go admin@example.com YourPassword123
```

Или через bash скрипт:
```bash
cd /opt/website-teacher/Backend
export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=ql1OdITTGh2023woLJToH3MaUGWYEMFH sslmode=disable"
./scripts/create-admin.sh admin@example.com YourPassword123
```

## Способ 2: Через Docker (если Go не установлен)

```bash
cd /opt/website-teacher

# Войдите в контейнер backend
docker-compose exec backend sh

# Внутри контейнера:
cd /root
export DATABASE_URL="host=postgres dbname=restapi_prod user=postgres password=ql1OdITTGh2023woLJToH3MaUGWYEMFH sslmode=disable"
go run scripts/create-admin.go admin@example.com YourPassword123
```

## Способ 3: Через SQL (требует генерации хеша)

Если у вас есть доступ к PostgreSQL и вы можете сгенерировать bcrypt хеш:

```bash
# Подключитесь к БД
docker-compose exec postgres psql -U postgres -d restapi_prod

# В psql выполните (замените email и хеш):
INSERT INTO admins (email, password_hash) 
VALUES ('admin@example.com', '$2a$10$хеш_пароля_здесь') 
ON CONFLICT (email) DO NOTHING;
```

## После создания админа

1. Откройте сайт: http://5.255.103.118:8000
2. Перейдите на страницу входа админа (обычно `/admin/login`)
3. Войдите с созданными данными

## Пример данных для админа

- **Email**: `admin@teacher.com` (или любой другой)
- **Password**: `Admin123456` (минимум 8 символов)

**Важно**: Используйте надежный пароль для продакшена!

