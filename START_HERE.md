# 🚀 Начните отсюда - Деплой на сервер

## Быстрый старт (3 шага)

### Шаг 1: Подключитесь к серверу

```bash
ssh root@5.255.103.118
# Пароль: 86aa1d6d890284c9e320ebe030641693
```

### Шаг 2: Выполните команды на сервере

Скопируйте и выполните все команды ниже на сервере:

```bash
# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh
systemctl start docker
systemctl enable docker

# Установка Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Клонирование репозитория
mkdir -p /opt/website-teacher
cd /opt/website-teacher
git clone https://github.com/1Dmitry2/website-teacher.git .

# Создание и настройка .env
cp env.example .env
nano .env
```

**В nano отредактируйте .env файл:**
- `POSTGRES_PASSWORD` - придумайте безопасный пароль
- `JWT_SECRET` - минимум 32 символа (можно сгенерировать: `openssl rand -base64 32`)
- `CORS_ORIGINS=http://5.255.103.118,http://yourdomain.com`
- `ADMIN_RESET_URL=http://5.255.103.118/admin/reset?token=`
- `USER_VERIFICATION_URL=http://5.255.103.118/verify-email?token=`
- `VITE_API_URL=http://5.255.103.118:8080`

Сохраните: `Ctrl+O`, `Enter`, `Ctrl+X`

### Шаг 3: Запустите контейнеры

```bash
cd /opt/website-teacher
docker-compose up -d --build
```

### Шаг 4: Запустите миграции

```bash
# Установите golang-migrate
wget https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz
tar -xzf migrate.linux-amd64.tar.gz
sudo mv migrate /usr/local/bin/
rm migrate.linux-amd64.tar.gz

# Запустите миграции (замените YOUR_PASSWORD на пароль из .env)
cd Backend
export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=YOUR_PASSWORD sslmode=disable"
./scripts/migrate.sh up
```

### Шаг 5: Откройте порты

```bash
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 8000/tcp # Nginx
ufw enable
```

## ✅ Готово!

Откройте в браузере: **http://5.255.103.118:8000**

## Полезные команды

```bash
# Просмотр логов
cd /opt/website-teacher
docker-compose logs -f

# Перезапуск
docker-compose restart

# Остановка
docker-compose down

# Обновление (после git pull)
docker-compose up -d --build
```

## Дополнительная документация

- [QUICK_DEPLOY.md](QUICK_DEPLOY.md) - Подробная инструкция
- [Deploy/SERVER_DEPLOY.md](Deploy/SERVER_DEPLOY.md) - Полная документация
- [Deploy/AUTO_DEPLOY.md](Deploy/AUTO_DEPLOY.md) - Автоматический деплой

