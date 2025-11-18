# Автоматический деплой

## Вариант 1: Копирование скрипта на сервер (рекомендуется)

1. **Скопируйте скрипт на сервер:**
```bash
scp Deploy/setup-server.sh root@5.255.103.118:/root/
```

2. **Подключитесь к серверу:**
```bash
ssh root@5.255.103.118
# Пароль: 86aa1d6d890284c9e320ebe030641693
```

3. **Запустите скрипт:**
```bash
chmod +x /root/setup-server.sh
/root/setup-server.sh
```

4. **Настройте .env:**
```bash
cd /opt/website-teacher
nano .env
```

5. **Запустите контейнеры:**
```bash
docker-compose up -d --build
```

## Вариант 2: Выполнение команд напрямую

Выполните эти команды на сервере:

```bash
# Подключитесь к серверу
ssh root@5.255.103.118

# Установите Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh
systemctl start docker
systemctl enable docker

# Установите Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Клонируйте репозиторий
cd /opt
git clone https://github.com/1Dmitry2/website-teacher.git
cd website-teacher

# Создайте .env
cp env.example .env
nano .env  # Отредактируйте настройки

# Запустите контейнеры
docker-compose up -d --build
```

## Настройка .env файла

Откройте `.env` и обязательно измените:

```bash
POSTGRES_PASSWORD=ваш_безопасный_пароль_здесь
JWT_SECRET=ваш_секретный_ключ_минимум_32_символа_длинный
CORS_ORIGINS=http://5.255.103.118,http://yourdomain.com
ADMIN_RESET_URL=http://5.255.103.118/admin/reset?token=
USER_VERIFICATION_URL=http://5.255.103.118/verify-email?token=
VITE_API_URL=http://5.255.103.118:8080
```

## Запуск миграций

После запуска контейнеров:

```bash
# Установите golang-migrate
wget https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz
tar -xzf migrate.linux-amd64.tar.gz
sudo mv migrate /usr/local/bin/
rm migrate.linux-amd64.tar.gz

# Запустите миграции
cd /opt/website-teacher/Backend
export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=ВАШ_ПАРОЛЬ_ИЗ_ENV sslmode=disable"
./scripts/migrate.sh up
```

## Открытие портов

```bash
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 8000/tcp # Nginx
ufw enable
```

## Проверка

Откройте в браузере: `http://5.255.103.118:8000`

Проверьте логи:
```bash
cd /opt/website-teacher
docker-compose logs -f
```

