#!/bin/bash

# Скрипт для настройки сервера
# Запустите этот скрипт на сервере после подключения по SSH

set -e

PROJECT_DIR="/opt/website-teacher"
GIT_REPO="https://github.com/1Dmitry2/website-teacher.git"

echo "🚀 Настройка сервера для деплоя..."

# Устанавливаем Docker
if ! command -v docker &> /dev/null; then
    echo "📥 Установка Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
    systemctl start docker
    systemctl enable docker
    echo "✅ Docker установлен"
else
    echo "✅ Docker уже установлен"
fi

# Устанавливаем Docker Compose
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null 2>&1; then
    echo "📥 Установка Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    echo "✅ Docker Compose установлен"
else
    echo "✅ Docker Compose уже установлен"
fi

# Создаем директорию проекта
mkdir -p $PROJECT_DIR
cd $PROJECT_DIR

# Клонируем или обновляем репозиторий
if [ -d ".git" ]; then
    echo "🔄 Обновление репозитория..."
    git fetch origin
    git reset --hard origin/main
    echo "✅ Репозиторий обновлен"
else
    echo "📥 Клонирование репозитория..."
    git clone $GIT_REPO .
    echo "✅ Репозиторий склонирован"
fi

# Создаем .env если его нет
if [ ! -f ".env" ]; then
    echo "📝 Создание .env файла из шаблона..."
    cp env.example .env
    echo ""
    echo "⚠️  ВАЖНО: Отредактируйте .env файл!"
    echo "   Выполните: nano .env"
    echo ""
    echo "Обязательно измените:"
    echo "  - POSTGRES_PASSWORD"
    echo "  - JWT_SECRET (минимум 32 символа)"
    echo "  - CORS_ORIGINS (укажите IP сервера: http://5.255.103.118)"
    echo "  - ADMIN_RESET_URL и USER_VERIFICATION_URL"
else
    echo "✅ .env файл уже существует"
fi

echo ""
echo "✅ Настройка завершена!"
echo ""
echo "📋 Следующие шаги:"
echo "1. Отредактируйте .env: nano .env"
echo "2. Запустите контейнеры: docker-compose up -d --build"
echo "3. Запустите миграции (см. QUICK_DEPLOY.md)"

