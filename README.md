## 📝 Ozon Habr API

Этот проект реализует API для озона на подобие хабра. Он позволет работать с постами и комментариями, используя GraphQL. Система поддерживает создание постов, добавление комментариев, управление включением/выключением комментариев для постов, а также подписку на добавление новых комментариев.

## 🚀 Запуск проекта
Чтобы запустить проект, выполните следующие шаги:

Убедитесь, что у вас установлены Docker и Docker Compose.

* Склонируйте репозиторий:
```
git clone https://github.com/nabishec/ozon_habr_api.git
cd ozon_habr_api
```
**В случае если вы хотите чтобsы хранение данных было in-memory то в Dockerfile надо изменить**

[`RUN go build -o main ./cmd/main.go`](Dockerfile#L7)

**на**

`RUN go build -o main ./cmd/main.go`

* Запустите проект с помощью Docker Compose:
```
docker-compose up
```
* После запуска откройте браузер и перейдите по адресу:
```
http://localhost:8080
```

## 📖 Документация API

Документация API доступна по адресу: 
* **http://localhost:8080/graphql**

Для тестирования API можно использовать любые GraphQL клиенты, например Insomnia или GraphiQL.

## 🛠 Архитектура и решения

В данном проекте используется структура данных для постов и комментариев с поддержкой пагинации, что позволяет эффективно работать с большими объемами данных. 
    
Кроме того:     
Комментарии хранятся в базе данных с использованием системы материализованных путей. Это позволяет эффективно строить иерархии комментариев и их ответов.
Redis кэширование используется для хранения комментариев на протяжении 15 минут, что ускоряет работу системы за счет сокращения количества запросов к базе данных.

## 📁 Структура проекта
<details>
    <summary style="display: inline-flex; align-items: center;">
        <b>Показать структуру </b>
    </summary>

`ozon_habr_api/`<br>
`├── cmd/`<br>
`│   ├── db_connection/`<br>
`│   │   ├──` [`cache.go`](./cmd/db_connection/cache.go)                (Подключение и настройка Redis для кэширования)<br>
`│   │   └──` [`database.go`](./cmd/db_connection/database.go)              (Подключение и настройка PostgreSQL)<br>
`│   ├── server/`<br>
`│   │   └──` [`server.go`](./cmd/server/server.go)               (Настройка и запуск GraphQL сервера)<br>
`│   └──` [`main.go`](./cmd/main.go)                     (Основная точка входа, настройка и запуск приложения)<br>
`├── graph/`<br>
`│   ├── model/`<br>
`│   │   └──` [`models_gen.go`](./graph/model/models_gen.go)           (Автоматически сгенерированные GraphQL модели)<br>
`│   ├──` [`generated.go`](./graph/generated.go)                 (Сгенерированный код GraphQL (gqlgen))<br>
`│   ├──` [`resolver.go`](./graph/resolver.go)                 (Основные резолверы GraphQL)<br>
`│   ├──` [`schema.graphqls`](./graph/schema.graphqls)             (Определение GraphQL схемы)<br>
`│   └──` [`schema.resolvers.go`](./graph/schema.resolvers.go)         (Реализация резолверов GraphQL)<br>
`├── internal/`<br>
`│   ├── handlers/`<br>
`│   │   ├── comment_query/`                (Обработчики логики запросов комментариев)<br>
`│   │   │   ├──` [`interface.go`](./internal/handlers/comment_query/interface.go)        (Интерфейс для запросов комментариев)<br>
`│   │   │   └──` [`query.go`](./internal/handlers/comment_query/query.go)        (Реализация запросов комментариев)<br>
`│   │   ├── post_mutation/`          (Обработчики логики мутаций постов)<br>
`│   │   │   ├──` [`interface.go`](./internal/handlers/post_mutation/interface.go)        (Интерфейс для мутаций постов)<br>
`│   │   │   └──` [`mutations.go`](./internal/handlers/post_mutation/mutations.go)        (Реализация мутаций постов)<br>
`│   │   └── post_query/`          (Обработчики логики запросов постов)<br>
`│   │       ├──` [`interface.go`](./internal/handlers/post_query/interface.go)        (Интерфейс для запросов постов)<br>
`│   │       └──` [`query.go`](./internal/handlers/post_query/query.go)        (Реализация запросов постов)<br>
`│   ├── lib/`<br>
`│   │   └── cursor/`<br>
`│   │       └──` [`cursor.go`](./internal/lib/cursor/cursor.go)        (Функции для работы с курсорами в пагинации)<br>
`│   ├── model/`<br>
`│   │   └──` [`model.go`](./internal/model/model.go)                (Внутренние модели данных)<br>
`│   └── storage/`<br>
`│       ├── db/` (Реализация хранилища данных в памяти) <br>
`│       │   └──` [`resolvers.go`](./internal/storage/db/resolvers.go)        (Реализация методов для работы с базой данных PostgreSQL)<br>
`│       ├── in-memory/` (Реализация хранилища данных в памяти) <br>
`│       │   └──` [`resolvers.go`](./internal/storage/in-memory/resolvers.go)        (Реализация методов для работы с даннми в памяти)<br>
`│       ├──` [`interface.go`](./internal/storage/interface.go)            (Интерфейс для хранилища данных (PostgreSQL, in-memory))<br>
`│       └──` [`storage_errors.go`](./internal/storage/storage_errors.go)     (Хранит ошибки бизнес логики)<br>
`├── migrations/`<br>
`│   └──` [`001_create_tables.up.sql`](./migrations/001_create_tables.up.sql)    (SQL скрипт для миграции базы данных (создание таблиц))<br>
`├── tools/`<br>
`│    └──` [`tools.go`](./tools/tools.go)                   (Инструменты для генерации кода gqlgen)<br>
`├──` [`.env`](./.env)                            (Файл с переменными окружения (настройки базы данных, Redis и т.д.))<br>
`├──` [`.gitignore`](./.gitignore)                      (Список игнорируемых файлов и директорий для Git)<br>
`├──` [`docker-compose.yml`](./docker-compose.yml)              (Конфигурация Docker Compose для запуска приложения и зависимостей)<br>
`├──` [`Dockerfile`](./Dockerfile)                      (Инструкции для сборки Docker образа)<br>
`├──` [`go.mod`](./go.mod)                          (Файл зависимостей Go)<br>
`├──` [`go.sum`](./go.sum)                          (Файл с контрольными суммами зависимостей Go)<br>
`├──` [`gqlgen.yml`](./gqlgen.yml)                      (Конфигурационный файл для gqlgen)<br>
`├──` [`LICENSE`](./LICENSE)                         (Лицензия проекта)<br>
`└──` [`README.md`](./README.md)                       (Файл с описанием проекта)<br>

</details>

## 📝 Структура данных
## 💡 Примечания
## 🧪 Тестирование
## 📜 Лицензия
