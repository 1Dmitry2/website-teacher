#!/bin/bash

# Скрипт для быстрого деплоя
# Скопируйте и выполните эти команды

echo "📋 Выполните эти команды на сервере:"
echo ""
echo "1. Подключитесь к серверу:"
echo "2. Выполните на сервере:"
echo ""
cat << 'SERVER_COMMANDS'
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

# Создание .env
cp env.example .env
nano .env  # Отредактируйте настройки!

# Запуск контейнеров
docker-compose up -d --build
SERVER_COMMANDS

echo ""
echo "3. После настройки .env запустите миграции (см. QUICK_DEPLOY.md)"

