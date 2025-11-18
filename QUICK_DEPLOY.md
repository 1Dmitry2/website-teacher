# Быстрый деплой на сервер

## Шаг 1: Закоммитьте и запушьте изменения

```bash
git add .
git commit -m "Add deployment configuration"
git push origin main
```

## Шаг 2: Подключитесь к серверу

```bash
ssh root@5.255.103.118
# Пароль: 86aa1d6d890284c9e320ebe030641693
```

## Шаг 3: Установите Docker и Docker Compose

```bash
# Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh
systemctl start docker
systemctl enable docker

# Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

## Шаг 4: Клонируйте репозиторий

**Если репозиторий публичный (HTTPS):**
```bash
cd /opt
git clone https://github.com/1Dmitry2/website-teacher.git
cd website-teacher
```

**Если репозиторий приватный (SSH):**
```bash
# Сначала настройте SSH ключ на сервере (см. ниже)
cd /opt
git clone git@github.com:1Dmitry2/website-teacher.git
cd website-teacher
```

## Шаг 5: Настройте .env файл

```bash
cp env.example .env
nano .env
```

**Обязательно измените:**
```
POSTGRES_PASSWORD=ваш_безопасный_пароль
JWT_SECRET=ваш_секретный_ключ_минимум_32_символа
CORS_ORIGINS=http://5.255.103.118,http://yourdomain.com
ADMIN_RESET_URL=http://5.255.103.118/admin/reset?token=
USER_VERIFICATION_URL=http://5.255.103.118/verify-email?token=
VITE_API_URL=http://5.255.103.118:8080
```

## Шаг 6: Запустите контейнеры

```bash
docker-compose up -d --build
```

## Шаг 7: Запустите миграции

```bash
# Установите golang-migrate
wget https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz
tar -xzf migrate.linux-amd64.tar.gz
sudo mv migrate /usr/local/bin/
rm migrate.linux-amd64.tar.gz

# Запустите миграции
cd Backend
export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=ВАШ_ПАРОЛЬ_ИЗ_ENV sslmode=disable"
./scripts/migrate.sh up
```

## Шаг 8: Откройте порты в файрволе

```bash
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 8000/tcp # Nginx
ufw enable
```

## Готово! 🎉

Откройте в браузере: `http://5.255.103.118:8000`

---

## Настройка SSH ключа для приватного репозитория

Если репозиторий приватный, нужно настроить SSH ключ:

```bash
# На сервере
ssh-keygen -t ed25519 -C "server@website-teacher"
cat ~/.ssh/id_ed25519.pub
```

Скопируйте публичный ключ и добавьте его в GitHub:
1. GitHub → Settings → SSH and GPG keys → New SSH key
2. Вставьте ключ

---

## Полезные команды

```bash
# Просмотр логов
docker-compose logs -f

# Перезапуск
docker-compose restart

# Остановка
docker-compose down

# Обновление (после git pull)
docker-compose up -d --build
```

Подробная документация: [Deploy/SERVER_DEPLOY.md](Deploy/SERVER_DEPLOY.md)

