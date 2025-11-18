#!/bin/bash

# Скрипт для деплоя на сервер
# Использование: ./deploy.sh

set -e

SERVER_IP="5.255.103.118"
SERVER_USER="root"
GIT_REPO="git@github.com:1Dmitry2/website-teacher.git"
PROJECT_DIR="/opt/website-teacher"

echo "🚀 Начинаем деплой на сервер $SERVER_IP..."

# Проверяем, что мы на правильной ветке
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "⚠️  Внимание: вы не на ветке main. Текущая ветка: $CURRENT_BRANCH"
    read -p "Продолжить? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Проверяем, что все изменения закоммичены
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  Есть незакоммиченные изменения!"
    echo "Пожалуйста, закоммитьте изменения перед деплоем:"
    echo "  git add ."
    echo "  git commit -m 'Prepare for deployment'"
    echo "  git push"
    exit 1
fi

echo "📦 Подключение к серверу и проверка Docker..."

# Подключаемся к серверу и выполняем команды
ssh $SERVER_USER@$SERVER_IP bash << EOF
    set -e
    
    # Проверяем наличие Docker
    if ! command -v docker &> /dev/null; then
        echo "📥 Установка Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sh get-docker.sh
        rm get-docker.sh
        systemctl start docker
        systemctl enable docker
    fi

    # Проверяем наличие Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null 2>&1; then
        echo "📥 Установка Docker Compose..."
        curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-\$(uname -s)-\$(uname -m)" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
    fi

    # Создаем директорию проекта
    mkdir -p $PROJECT_DIR

    # Переходим в директорию проекта
    cd $PROJECT_DIR

    # Клонируем или обновляем репозиторий
    if [ -d ".git" ]; then
        echo "🔄 Обновление репозитория..."
        git fetch origin
        git reset --hard origin/main
    else
        echo "📥 Клонирование репозитория..."
        git clone $GIT_REPO .
    fi

    # Проверяем наличие .env файла
    if [ ! -f ".env" ]; then
        echo "⚠️  Файл .env не найден!"
        echo "Создайте .env из env.example и настройте его:"
        echo "  cp env.example .env"
        echo "  nano .env"
        exit 1
    fi

    echo "🐳 Остановка старых контейнеров..."
    docker-compose down || true

    echo "🔨 Сборка и запуск контейнеров..."
    docker-compose up -d --build

    echo "✅ Деплой завершен!"
    echo "Проверьте статус: docker-compose ps"
    echo "Просмотр логов: docker-compose logs -f"
EOF

echo ""
echo "✅ Деплой завершен!"
echo ""
echo "📋 Следующие шаги:"
echo "1. Подключитесь к серверу: ssh $SERVER_USER@$SERVER_IP"
echo "2. Перейдите в директорию: cd $PROJECT_DIR"
echo "3. Настройте .env файл (если еще не настроен): nano .env"
echo "4. Запустите миграции (см. Deploy/README.md)"
echo "5. Проверьте логи: docker-compose logs -f"

