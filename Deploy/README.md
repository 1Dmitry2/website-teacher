# Деплой сайта-визитки для препода

Этот проект готов к деплою с использованием Docker и Docker Compose.

## Структура проекта

- **Backend**: Go API сервер (Gin framework)
- **Frontend**: Vue 3 приложение (Vite)
- **Database**: PostgreSQL
- **Nginx**: Reverse proxy для продакшена

## Быстрый старт

### 1. Подготовка

1. Скопируйте `env.example` в `.env`:
```bash
cp env.example .env
```

2. Отредактируйте `.env` файл и укажите:
   - Пароль для PostgreSQL
   - JWT секрет (минимум 32 символа)
   - SMTP настройки для отправки email
   - URL вашего домена для ссылок в письмах
   - CORS origins (разделенные запятыми)

### 2. Запуск

```bash
# Собрать и запустить все сервисы
docker-compose up -d

# Посмотреть логи
docker-compose logs -f

# Остановить
docker-compose down
```

### 3. Миграции базы данных

Для запуска миграций используйте golang-migrate:

```bash
# Установите golang-migrate
# macOS: brew install golang-migrate
# Linux: https://github.com/golang-migrate/migrate

# Запустите миграции
cd Backend
export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=YOUR_PASSWORD sslmode=disable"
./scripts/migrate.sh up
```

Или используйте Docker:

```bash
docker-compose exec backend sh -c "cd /root && migrate -path ./migrations -database 'postgres://postgres:YOUR_PASSWORD@postgres:5432/restapi_prod?sslmode=disable' up"
```

## Конфигурация

### Переменные окружения

Все настройки можно изменить через `.env` файл:

- `POSTGRES_*` - настройки базы данных
- `JWT_SECRET` - секретный ключ для JWT токенов
- `SMTP_*` - настройки SMTP сервера для отправки email
- `CORS_ORIGINS` - разрешенные источники для CORS (через запятую)
- `ADMIN_RESET_URL` - URL для сброса пароля админа
- `USER_VERIFICATION_URL` - URL для верификации email пользователей

### Порты

По умолчанию:
- Frontend: `80` (внутри контейнера)
- Backend: `8080` (внутри контейнера)
- PostgreSQL: `5432` (внутри контейнера)
- Nginx: `8000` (внешний порт)

Измените порты в `.env` файле при необходимости.

## Структура сервисов

```
┌─────────┐
│  Nginx  │ :8000 (внешний порт)
└────┬────┘
     │
     ├──> Frontend :80
     ├──> Backend :8080
     └──> PostgreSQL :5432
```

## Email настройка

Для работы верификации email и сброса пароля нужно настроить SMTP.

### Gmail

1. Включите двухфакторную аутентификацию
2. Создайте пароль приложения: https://myaccount.google.com/apppasswords
3. Укажите в `.env`:
   ```
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your-email@gmail.com
   SMTP_PASSWORD=your-app-password
   SMTP_FROM=your-email@gmail.com
   ```

### Другие провайдеры

Используйте настройки вашего SMTP провайдера.

## Деплой на продакшен

### 1. Подготовка сервера

- Установите Docker и Docker Compose
- Настройте домен и DNS
- Настройте SSL сертификат (Let's Encrypt)

### 2. Обновление конфигурации

1. Обновите `.env` с продакшен значениями
2. Установите правильные URL для email ссылок:
   ```
   ADMIN_RESET_URL=https://yourdomain.com/admin/reset?token=
   USER_VERIFICATION_URL=https://yourdomain.com/verify-email?token=
   ```
3. Настройте CORS:
   ```
   CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
   ```

### 3. SSL/HTTPS

Для продакшена рекомендуется использовать Nginx с SSL. Обновите `Deploy/nginx.conf` для поддержки HTTPS.

### 4. Бэкапы

Настройте регулярные бэкапы базы данных:
```bash
docker-compose exec postgres pg_dump -U postgres restapi_prod > backup.sql
```

## Мониторинг

Просмотр логов:
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

## Обновление

```bash
# Остановить
docker-compose down

# Обновить код
git pull

# Пересобрать и запустить
docker-compose up -d --build
```

## Troubleshooting

### База данных не подключается

Проверьте:
- Правильность `DATABASE_URL` в `.env`
- Что PostgreSQL контейнер запущен: `docker-compose ps`
- Логи PostgreSQL: `docker-compose logs postgres`

### CORS ошибки

Убедитесь, что `CORS_ORIGINS` в `.env` содержит ваш домен.

### Email не отправляется

- Проверьте SMTP настройки в `.env`
- Проверьте логи backend: `docker-compose logs backend`
- Приложение будет работать без SMTP, но письма не будут отправляться

## Полезные команды

```bash
# Перезапустить сервис
docker-compose restart backend

# Войти в контейнер
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d restapi_prod

# Очистить все (включая volumes)
docker-compose down -v

# Просмотр использования ресурсов
docker stats
```

