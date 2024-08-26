# notes_api

Часть сервиса, предоставляющего REST API интерфейс с методами:
* Добавление заметки
* Вывод списка заметок

### Используемый стек

* **Golang 1.22**
* **Echo** основной веб фреймворк
* **PostgreSQL** как основная БД
* **golang-jwt/jwt** для jwt
* **swaggo/swag** swagger документация API
* **golang-migrate/migrate** для миграций бд
* **logrus** для логирования
* **testify** для тестирования
* **Docker и Docker Compose** для быстрого развертывания


### Вопросы по тестовому заданию

В процессе разработки я столкнулся с некоторыми вопросами касательно некоторых моментов:  

**Валидация орфографических ошибок**  
В задании указано, что при сохранении заметок необходимо проверять текст на ошибки.
Появился вопрос, каким образом это реализовывать: при сохранении или отдельный маршрут.
Решил **_сделать отдельным маршрутом_** валидацию, а не проверять текст на ошибки перед сохранением, потому что так будет логически правильно
(будет странно, если мы не будем сохранять заметку только потому что там есть орфографические ошибки).


### Prerequisites
- Docker, Docker Compose installed

### Getting started

* Добавить репозиторий к себе
* Создать .env файл в директории с проектом и заполнить информацией из .env.example
* Применить `make keys` (или самостоятельно создать приватный и публичный rsa ключи)

### Usage

Запустить сервис можно с помощью `make compose-up` (или `docker-compose up -d --build`)

Документация доступна по адресу `http://localhost:8080/swagger/`

Запуск тестов доступен с помощью команды `make tests`


### Примеры запросов

* [Регистрация](#регистрация)
* [Аутентификация](#вход-аутентификация)
* [Создание заметки](#создание-заметки)
* [Получение списка заметок](#получение-списка-заметок)
* [Валидирование текста](#валидирование-текста-на-наличие-грамматических-ошибок)


#### Регистрация
```shell
curl -X 'POST' \
      'http://localhost:8080/auth/sign-up' \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{"username": "maks", "password": "abc"}'
```
Пример ответа:  
`201`

#### Вход, аутентификация
```shell
curl -X 'POST' \
      'http://localhost:8080/auth/sign-in' \
      -H 'accept: application/json' \
      -H 'Content-Type: application/json' \
      -d '{"password": "abc", "username": "maks"}'
```
Пример ответа:  
```json
{
  "token": "jwt-token"
}
```

#### Создание заметки
```shell
curl -X 'POST' \
      'http://localhost:8080/api/v1/notes/create' \
      -H 'accept: application/json' \
      -H 'Authorization: Bearer jwt-token' \
      -H 'Content-Type: application/json' \
      -d '{"text": "foobar", "title": "Hello world"}'
```
Пример ответа:
```json
{
  "note_id": 1
}
```

#### Получение списка заметок
```shell
curl -X 'GET' \
      'http://localhost:8080/api/v1/notes/list' \
      -H 'accept: application/json' \
      -H 'Authorization: Bearer jwt-token' \
      -H 'Content-Type: application/json' \
      -d '{"limit": 20, "offset": 0, "sort": "id"}'
```
Пример ответа:
```json
[
  {
    "id": 1,
    "title": "Hello world",
    "text": "foobar",
    "created_at": "2024-08-26T06:39:31.989361Z"
  }
]
```

#### Валидирование текста на наличие грамматических ошибок
```shell
curl -X 'POST' \
      'http://localhost:8080/api/v1/notes/validate' \
      -H 'accept: application/json' \
      -H 'Authorization: Bearer jwt-token' \
      -H 'Content-Type: application/json' \
      -d '{"text": "Превет, как дела?"}'
```
Пример ответа:
```json
[
  {
    "type": "Unknown word",
    "position": 0,
    "row": 0,
    "column": 0,
    "length": 6,
    "word": "Превет",
    "replacements": [
      "Привет",
      "Превед",
      "Приветь"
    ]
  }
]
```


### Тестовое задание
Необходимо спроектировать и реализовать на Golang сервис, предоставляющий REST API интерфейс с методами:
* Добавление заметки
* Вывод списка заметок

При сохранении заметок необходимо орфографические ошибки валидировать при помощи сервиса [Яндекс.Спеллер](https://yandex.ru/dev/speller/).  
Также необходимо реализовать аутентификацию и авторизацию. Пользователи должны иметь доступ только к своим заметкам.