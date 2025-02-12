# URL-SHORTENING-SERVICE

## Описание
Проект представляет собой систему для укорочения URL-адресов. Клиент отправляет оригинальный URL адрес и в ответ получает сокращённый(alias). Так же по сокращенному адресу он может получить оригинальный(url). Пользователь может выбрать тип хранилища URL адресов(postgres/in-memory), указав в соответствующем параметре(STORAGE) в `.env` файле. Сервис работает на `gRPC`.

### Конфигурация

Вы можете настроить ваше окружение с помощью `.env` файла. Также предлагаем воспользоваться рекомендуемым файлом конфигурации:

```env
# APP Configuration
ENV=local
LOG_LEVEL=debug

# gRPC Server Configuration
PORT=44044
TIMEOUT=10h

# storage Configuration
STORAGE=postgres
DB_REPLICAS=1
MIGRATIONS_PATH=file://migrations/

# PostgreSQL Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DBNAME=url_shortening
POSTGRES_SSLMODE=disable
```

### Установка

Проект запускается по команде:

`make compose`

Для выбора хранилища в файле `.env` можно изменить переменные `STORAGE` и `DB_REPLICAS`:

`STORAGE="postgres"/"in-memory"`

`DB_REPLICAS=0` (in-memory)

`DB_REPLICAS=1` (postgres)

### API Endpoints

#### 1. Сохранение оригинального URL и возврат сокращённого:

**Method:** `SaveUrl`
- **Request parameters:**

    - `url` - Оригинальный URL

- **Response:**

    - `alias` - Сокращённый URL

#### 2. Получение оригинального URL по сокращённому:

**Method:** `RedirectUrl`

- **Request parameters:**

    - `alias` - Сокращённый URL

- **Response:**

    - `url` - Оригинальный URL

### Запуск тестов по команде:

`make test`
