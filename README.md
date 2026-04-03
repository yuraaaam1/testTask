# Subscriptions API

REST-сервис для агрегации данных об онлайн подписках пользователей.

## Технологии

- Go 1.25
- PostgreSQL 18
- Docker / Docker Compose
- Gin — HTTP роутер
- golang-migrate — миграции БД
- zap — логирование
- Swagger — документация API

## Запуск

1. Склонировать репозиторий:
```bash
git clone https://github.com/yuraaaam1/testTask.git
cd testTask
Создать .env файл на основе .env.example:

cp .env.example .env
Заполнить .env своими значениями

Запустить через Docker Compose:


docker compose up --build
API
Swagger документация доступна по адресу:


http://localhost:8080/swagger/index.html
Эндпоинты
Метод	URL	Описание
POST	/api/v1/subscriptions	Создать подписку
GET	/api/v1/subscriptions	Получить список подписок
GET	/api/v1/subscriptions/:id	Получить подписку по ID
PUT	/api/v1/subscriptions/:id	Обновить подписку
DELETE	/api/v1/subscriptions/:id	Удалить подписку
GET	/api/v1/subscriptions/cost	Подсчёт суммарной стоимости
Пример запроса на создание подписки

{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
}
Подсчёт стоимости

GET /api/v1/subscriptions/cost?user_id=<uuid>&date_from=01-2025&date_to=12-2025
Параметры:

date_from — начало периода (обязательный, формат MM-YYYY)
date_to — конец периода (обязательный, формат MM-YYYY)
user_id — фильтр по пользователю (опциональный)
service_name — фильтр по названию сервиса (опциональный)

