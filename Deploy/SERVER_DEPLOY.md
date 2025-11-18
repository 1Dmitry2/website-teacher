# Инструкция по деплою на сервер

## Данные сервера
- **IP адрес**: 5.255.103.118
- **Пользователь**: root
- **Пароль**: (указан отдельно)

## Вариант 1: Автоматический деплой (рекомендуется)

### Шаг 1: Закоммитьте изменения

```bash
git add .
git commit -m "Add deployment configuration"
git push origin main
```

### Шаг 2: Запустите скрипт деплоя

```bash
chmod +x Deploy/deploy.sh
./Deploy/deploy.sh
```

Скрипт автоматически:
- Установит Docker и Docker Compose (если нужно)
- Клонирует/обновит репозиторий на сервере
- Соберет и запустит контейнеры

## Вариант 2: Ручной деплой

### Шаг 1: Подключитесь к серверу

```bash
ssh root@5.255.103.118
# Введите пароль при запросе
```

### Шаг 2: Установите Docker и Docker Compose

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

# Проверка
docker --version
docker-compose --version
```

### Шаг 3: Клонируйте репозиторий

```bash
cd /opt
git clone git@github.com:1Dmitry2/website-teacher.git
cd website-teacher
```

**Или если используете HTTPS:**
```bash
git clone https://github.com/1Dmitry2/website-teacher.git
```

### Шаг 4: Настройте переменные окружения

```bash
cp env.example .env
nano .env
```

**Важно настроить:**
- `POSTGRES_PASSWORD` - безопасный пароль
- `JWT_SECRET` - минимум 32 символа
- `CORS_ORIGINS` - укажите IP сервера или домен: `http://5.255.103.118,http://yourdomain.com`
- `ADMIN_RESET_URL` - `http://5.255.103.118/admin/reset?token=` (или ваш домен)
- `USER_VERIFICATION_URL` - `http://5.255.103.118/verify-email?token=` (или ваш домен)
- `VITE_API_URL` - `http://5.255.103.118:8080` (или через домен)
- SMTP настройки (если нужна отправка email)

### Шаг 5: Запустите контейнеры

```bash
docker-compose up -d --build
```

### Шаг 6: Запустите миграции базы данных

```bash
# Установите golang-migrate на сервере
# Для Ubuntu/Debian:
wget https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz
tar -xzf migrate.linux-amd64.tar.gz
sudo mv migrate /usr/local/bin/
rm migrate.linux-amd64.tar.gz

# Запустите миграции
cd Backend
export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=YOUR_PASSWORD sslmode=disable"
./scripts/migrate.sh up
```

**Или через Docker:**
```bash
docker-compose exec backend sh -c "cd /root && migrate -path ./migrations -database 'postgres://postgres:YOUR_PASSWORD@postgres:5432/restapi_prod?sslmode=disable' up"
```

### Шаг 7: Проверьте статус

```bash
docker-compose ps
docker-compose logs -f
```

## Настройка файрвола

Откройте необходимые порты:

```bash
# Если используется ufw
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS (если настроен)
ufw allow 8000/tcp  # Nginx (если используете этот порт)
ufw enable

# Или если используется firewalld
firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --permanent --add-port=443/tcp
firewall-cmd --permanent --add-port=8000/tcp
firewall-cmd --reload
```

## Настройка домена (опционально)

Если у вас есть домен:

1. Настройте DNS записи, чтобы домен указывал на IP `5.255.103.118`
2. Обновите `.env` файл:
   - `CORS_ORIGINS=https://yourdomain.com,https://www.yourdomain.com`
   - `ADMIN_RESET_URL=https://yourdomain.com/admin/reset?token=`
   - `USER_VERIFICATION_URL=https://yourdomain.com/verify-email?token=`
3. Настройте SSL сертификат (Let's Encrypt) через Nginx

## Обновление приложения

```bash
ssh root@5.255.103.118
cd /opt/website-teacher
git pull origin main
docker-compose up -d --build
```

## Полезные команды

```bash
# Просмотр логов
docker-compose logs -f
docker-compose logs -f backend
docker-compose logs -f frontend

# Перезапуск сервиса
docker-compose restart backend

# Остановка всех сервисов
docker-compose down

# Остановка с удалением volumes (ОСТОРОЖНО - удалит данные БД!)
docker-compose down -v

# Вход в контейнер
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d restapi_prod
```

## Troubleshooting

### Не могу подключиться по SSH
- Проверьте, что порт 22 открыт в файрволе
- Убедитесь, что SSH сервис запущен на сервере

### Контейнеры не запускаются
```bash
docker-compose logs
# Проверьте логи конкретного сервиса
```

### База данных не подключается
- Проверьте `DATABASE_URL` в `.env`
- Убедитесь, что PostgreSQL контейнер запущен: `docker-compose ps`

### CORS ошибки
- Убедитесь, что `CORS_ORIGINS` в `.env` содержит ваш домен/IP

## Безопасность

⚠️ **Важно для продакшена:**

1. Измените пароль root или создайте нового пользователя с sudo правами
2. Настройте SSH ключи вместо пароля
3. Используйте сильные пароли в `.env`
4. Настройте SSL/HTTPS
5. Регулярно обновляйте систему и Docker
6. Настройте бэкапы базы данных

