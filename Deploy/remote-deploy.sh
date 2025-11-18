#!/bin/bash

# Скрипт для деплоя на удаленный сервер
# Использование: ./remote-deploy.sh

set -e

SERVER_IP="5.255.103.118"
SERVER_USER="root"
SERVER_PASS="86aa1d6d890284c9e320ebe030641693"
GIT_REPO="https://github.com/1Dmitry2/website-teacher.git"
PROJECT_DIR="/opt/website-teacher"

echo "🚀 Начинаем деплой на сервер $SERVER_IP..."

# Функция для выполнения команд на сервере
run_remote() {
    sshpass -p "$SERVER_PASS" ssh -o StrictHostKeyChecking=no $SERVER_USER@$SERVER_IP "$@"
}

# Проверяем наличие sshpass
if ! command -v sshpass &> /dev/null; then
    echo "📦 Установка sshpass..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install hudochenkov/sshpass/sshpass
        else
            echo "❌ Нужно установить sshpass. Установите Homebrew и выполните: brew install hudochenkov/sshpass/sshpass"
            exit 1
        fi
    else
        # Linux
        sudo apt-get update && sudo apt-get install -y sshpass || sudo yum install -y sshpass
    fi
fi

echo "📦 Подключение к серверу..."

# Устанавливаем Docker если нужно
run_remote bash << 'REMOTE_SCRIPT'
    set -e
    
    # Проверяем наличие Docker
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

    # Проверяем наличие Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null 2>&1; then
        echo "📥 Установка Docker Compose..."
        curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
        echo "✅ Docker Compose установлен"
    else
        echo "✅ Docker Compose уже установлен"
    fi
REMOTE_SCRIPT

# Создаем директорию и клонируем/обновляем репозиторий
run_remote bash << REMOTE_SCRIPT
    set -e
    
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
        echo "📝 Создание .env файла..."
        cp env.example .env
        echo "⚠️  ВАЖНО: Отредактируйте .env файл перед запуском!"
        echo "   Выполните: ssh $SERVER_USER@$SERVER_IP 'cd $PROJECT_DIR && nano .env'"
    else
        echo "✅ .env файл уже существует"
    fi
REMOTE_SCRIPT

echo ""
echo "✅ Базовая настройка завершена!"
echo ""
echo "📋 Следующие шаги:"
echo ""
echo "1. Настройте .env файл на сервере:"
echo "   ssh $SERVER_USER@$SERVER_IP"
echo "   cd $PROJECT_DIR"
echo "   nano .env"
echo ""
echo "2. После настройки .env запустите контейнеры:"
echo "   cd $PROJECT_DIR"
echo "   docker-compose up -d --build"
echo ""
echo "3. Запустите миграции (см. QUICK_DEPLOY.md)"
echo ""

