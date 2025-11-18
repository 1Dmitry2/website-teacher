# Инструкция по деплою

## Быстрый старт

1. **Скопируйте файл с переменными окружения:**
   ```bash
   cp env.example .env
   ```

2. **Отредактируйте `.env` файл:**
   - Установите безопасный пароль для PostgreSQL
   - Установите JWT_SECRET (минимум 32 символа)
   - Настройте SMTP для отправки email (опционально)
   - Укажите ваши домены в CORS_ORIGINS, ADMIN_RESET_URL, USER_VERIFICATION_URL

3. **Запустите проект:**
   ```bash
   docker-compose up -d
   ```

4. **Запустите миграции базы данных:**
   ```bash
   # Установите golang-migrate: https://github.com/golang-migrate/migrate
   cd Backend
   export DATABASE_URL="host=localhost dbname=restapi_prod user=postgres password=YOUR_PASSWORD sslmode=disable"
   ./scripts/migrate.sh up
   ```

5. **Откройте в браузере:**
   - Frontend: http://localhost:80 (или порт из .env)
   - Через Nginx: http://localhost:8000 (или порт из .env)

## Подробная документация

См. [Deploy/README.md](Deploy/README.md) для полной документации.

## Важные замечания

- **Безопасность**: Обязательно измените все пароли и секреты в `.env` перед деплоем на продакшен
- **Email**: Настройка SMTP опциональна, но необходима для верификации email и сброса пароля
- **Миграции**: Не забудьте запустить миграции перед первым запуском
- **SSL**: Для продакшена рекомендуется настроить HTTPS через Nginx

