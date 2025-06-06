# 📝 Ozon Habr API
Проект реализует API для озона на подобие хабра. Он позволяет работать с постами и комментариями, используя GraphQL. Система поддерживает создание постов, добавление комментариев, управление включением/выключением комментариев для постов, а также подписку на добавление новых комментариев.

## 🚀 Запуск проекта

Чтобы запустить проект, выполните следующие шаги:

### Предварительные требования

- Docker и Docker Compose установлены на вашей системе
- Git для клонирования репозитория

### Шаги по установке

1. **Клонирование репозитория**
   ```bash
   git clone https://github.com/nabishec/ozon_habr_api.git
   cd ozon_habr_api
   ```

2. **Запуск в Docker-контейнерах**
   ```bash
   docker-compose up
   ```

3. **Использование API**
   
   После запуска откройте браузер и перейдите по адресу:
   ```
   http://localhost:8080
   ```

### Выбор хранилища данных

По умолчанию проект использует PostgreSQL для хранения данных. Если вы хотите использовать in-memory хранилище для тестирования, измените в Dockerfile строку:
```
CMD ["./main"]
```
и добавьте флаг [`"-s", "m"`](Dockerfile#L11):
```
CMD ["./main","-s", "m"]
```

## 📖 Документация API

Интерактивная GraphQL-playground консоль доступна по адресу:
* **http://localhost:8080**

### Примеры GraphQL запросов

<details>
    <summary><b>Получение списка всех постов</b></summary>
    
    query{
        posts{
            id
            title 
            text
            authorID
            commentsEnabled
            createDate
        }
    }

</details>

<details>
    <summary><b>Получение конкретного поста с комментариями</b></summary>
    
    query {
        post(postID: 1) {
            id
            title
            text
            comments(first: 5) {
                edges {
                    node {
                        id
                        text
                        authorID
                        createDate
                    }
                    cursor
                }
                pageInfo {
                    hasNextPage
                    endCursor
                }
            }
        }
    }

</details>

<details>
    <summary><b>Создание нового поста</b></summary>
    
    mutation {
        addPost(postInput: {
            authorID: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
            title: "Новый пост"
            text: "Содержимое поста"
            commentsEnabled: true
        }) {
            id
            title
            createDate
        }
    }

</details>

<details>
    <summary><b>Создание нового комментария</b></summary>
    
    mutation {
        addComment(commentInput: {
            authorID: "123e4567-e89b-12d3-a456-426614174000",
            postID: 1,
            parentID: 1, # ID существующего комментария
            text: "Это ответ на комментарий 1"
        }) {
            id
            text
            parentID
        }   
    }

</details>


<details>
    <summary><b>Подписка на новые комментарии</b></summary>
    
    subscription {
        commentAdded(postID: 1) {
            id
            text
            authorID
            createDate
        }
    }

</details>

Для тестирования API можно использовать любые GraphQL клиенты, например Insomnia, Postman или GraphiQL.

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
`│   ├──` [`schema.resolvers.go`](./graph/schema.resolvers.go)         (Реализация резолверов GraphQL)<br>
`│   └──` [`subscription.go`](./graph/subscription.go)         (Реализация структур и методов для управления подписками)<br>
`├── internal/`<br>
`│   ├── handlers/`<br>
`│   │   ├── comment_mutation/`                (Обработчики логики мутаций комментариев)<br>
`│   │   │   ├──` [`interface.go`](./internal/handlers/comment_mutation/interface.go)        (Интерфейс для мутаций комментариев)<br>
`│   │   │   ├──` [`mutations.go`](./internal/handlers/comment_mutation/mutations.go)        (Реализация мутаций комментариев)<br>
`│   │   │   ├──` [`mutation_test.go`](./internal/handlers/comment_mutation/mutation_test.go)     (Тесты для мутаций комментариев)<br>
`│   │   │   └──` [`comment_mutation_imp_mock_test.go`](./internal/handlers/comment_mutation/comment_mutation_imp_mock_test.go)  (Моки для тестирования мутаций комментариев)<br>
`│   │   ├── comment_query/`                (Обработчики логики запросов комментариев)<br>
`│   │   │   ├──` [`interface.go`](./internal/handlers/comment_query/interface.go)        (Интерфейс для запросов комментариев)<br>
`│   │   │   ├──` [`query.go`](./internal/handlers/comment_query/query.go)        (Реализация запросов комментариев)<br>
`│   │   │   ├──` [`query_test.go`](./internal/handlers/comment_query/query_test.go)     (Тесты для запросов комментариев)<br>
`│   │   │   └──` [`comment_query_imp_mock_test.go`](./internal/handlers/comment_query/comment_query_imp_mock_test.go)  (Моки для тестирования запросов комментариев)<br>
`│   │   ├── post_mutation/`          (Обработчики логики мутаций постов)<br>
`│   │   │   ├──` [`interface.go`](./internal/handlers/post_mutation/interface.go)        (Интерфейс для мутаций постов)<br>
`│   │   │   ├──` [`mutations.go`](./internal/handlers/post_mutation/mutations.go)        (Реализация мутаций постов)<br>
`│   │   │   ├──` [`mutations_test.go`](./internal/handlers/post_mutation/mutations_test.go)     (Тесты для мутаций постов)<br>
`│   │   │   └──` [`post_mut_imp_mock_test.go`](./internal/handlers/post_mutation/post_mut_imp_mock_test.go)  (Моки для тестирования мутаций постов)<br>
`│   │   └── post_query/`          (Обработчики логики запросов постов)<br>
`│   │       ├──` [`interface.go`](./internal/handlers/post_query/interface.go)        (Интерфейс для запросов постов)<br>
`│   │       ├──` [`query.go`](./internal/handlers/post_query/query.go)        (Реализация запросов постов)<br>
`│   │       ├──` [`query_test.go`](./internal/handlers/post_query/query_test.go)     (Тесты для запросов постов)<br>
`│   │       └──` [`post_query_imp_mock_test.go`](./internal/handlers/post_query/post_query_imp_mock_test.go)  (Моки для тестирования запросов постов)<br>
`│   ├── pkg/`<br>
`│   │   ├── cursor/`<br>
`│   │   |   └──` [`cursor.go`](./internal/pkg/cursor/cursor.go)        (Функции для работы с курсорами в пагинации)<br>
`│   │   └── errs/`<br>
`│   │       └──` [`errors.go`](./internal/pkg/errs/errors.go)        (Хранит ошибки бизнес логики)<br>
`│   ├── model/`<br>
`│   │   └──` [`model.go`](./internal/model/model.go)                (Внутренние модели данных)<br>
`│   └── storage/`<br>
`│       ├── db/` (Реализация хранилища данных в PostgreSQL) <br>
`│       │   └──` [`resolvers.go`](./internal/storage/db/resolvers.go)        (Реализация методов для работы с базой данных PostgreSQL)<br>
`│       ├── in-memory/` (Реализация хранилища данных в памяти) <br>
`│       │   └──` [`resolvers.go`](./internal/storage/in-memory/resolvers.go)        (Реализация методов для работы с данными в памяти)<br>
`│       └──` [`interface.go`](./internal/storage/interface.go)            (Интерфейс для хранилища данных (PostgreSQL, in-memory))<br>
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
`├──` [`cover.out`](./cover.out)                      (Файл с отчетом о покрытии кода тестами)<br>
`├──` [`image.png`](./image.png)                      (Изображение со статистикой покрытия кода тестами)<br>
`├──` [`LICENSE`](./LICENSE)                         (Лицензия проекта)<br>
`└──` [`README.md`](./README.md)                       (Файл с описанием проекта)<br>

</details>

## 📝 Структура данных

Проект использует эффективную иерархическую структуру данных для организации комментариев:
<details>
    <summary style="display: inline-flex; align-items: center;">
        <b>Показать структуру </b>
    </summary>

### Посты (Posts)
- **ID**: Уникальный идентификатор поста (BIGSERIAL)
- **AuthorID**: UUID автора поста
- **Title**: Заголовок поста
- **Text**: Содержимое поста
- **CommentsEnabled**: Флаг, указывающий, разрешены ли комментарии к посту
- **CreateDate**: Дата и время создания поста

### Комментарии (Comments)
- **ID**: Уникальный идентификатор комментария (BIGSERIAL)
- **AuthorID**: UUID автора комментария
- **PostID**: Идентификатор поста, к которому относится комментарий
- **ParentID**: Идентификатор родительского комментария (для вложенных комментариев)
- **Path**: Материализованный путь в формате LTREE для эффективного поиска и построения иерархии
- **Text**: Текст комментария
- **CreateDate**: Дата и время создания комментария

</details>

## 🛠 Особенности реализации:

<details>
    <summary><b>Показать</b></summary>

### 🐘 Работа с данными в PostgreSQL:
- **Материализованные пути**: Использование PostgreSQL LTREE для хранения иерархии комментариев обеспечивает высокую производительность при запросах вложенных структур
- **Оптимизированные индексы**: Созданы индексы по полям path, create_date и post_id для ускорения запросов
- **Эффективная организация комментариев**: Иерархическая структура комментариев с возможностью глубокой вложенности до любого уровня

### 🔄 Кэширование с использованием Redis:
- **Двухуровневое кэширование**: Использование Redis для кэширования запрашиваемых комментариев и веток обсуждений на 30 минут
- **Инвалидация кэша**: Автоматическое обновление кэша при создании новых комментариев для обеспечения актуальности данных

P.s. Если в кэше не находит ветки, то делается запрос в бд.

### 📊 Оптимизация API и производительности:
- **Пагинация на основе курсоров**: Эффективная пагинация результатов с сохранением контекста для больших наборов данных
- **GraphQL оптимизация**: Возможность запрашивать только необходимые поля и эффективная организация связанных данных
- **Параллельная обработка запросов**: Использование горутин для обработки тяжелых задач без блокировки основного потока

### 📡 Система подписок (WebSockets):
- **Паттерн Publish-Subscribe**: Реализация Pub/Sub для уведомлений о новых комментариях в реальном времени, где компоненты взаимодействуют через центральный механизм каналов
- **Потокобезопасное управление подписками**: Использование мьютексов для безопасного доступа к списку подписчиков в конкурентной среде
- **Автоматическая очистка ресурсов**: Корректное закрытие каналов и удаление неактивных подписчиков для предотвращения утечек памяти
- **Асинхронность**: Используются неблокирующие Go-каналы для передачи данных
- **Устойчивость к ошибкам**: Защита от паник при отправке данных в закрытые каналы с использованием отложенных функций
- **Масштабируемость**: Возможность подписки на события по конкретному идентификатору поста, что обеспечивает точечную доставку уведомлений

> **Примечание о паттернах**: В отличие от классического паттерна Observer, где наблюдатели напрямую регистрируются у наблюдаемого объекта, в данном проекте реализован паттерн Publish-Subscribe, который вводит промежуточный слой (брокер сообщений) между издателями и подписчиками. Это обеспечивает более высокую степень декомпозиции: издатели не знают о конкретных подписчиках, а подписчики не знают об издателях. Подписки группируются по идентификатору поста, что позволяет реализовать фильтрацию событий на уровне брокера.

</details>

## 🐳 Контейнеризация

Проект полностью готов к запуску в контейнерах:

- **Dockerfile**: Оптимизированный многостадийный образ для минимального размера и максимальной производительности
- **Docker Compose**: Полная конфигурация для запуска всей инфраструктуры (API, PostgreSQL, Redis) одной командой
- **Переменные окружения**: Настройка всех компонентов через переменные окружения для гибкого развертывания


Это обеспечивает:
- Идентичность сред разработки и производства
- Простоту горизонтального масштабирования
- Минимальное время развертывания
- Легкую интеграцию в любую облачную платформу (AWS, GCP, Azure)

## 💭 Примечания к архитектуре и инженерным решениям
Данный проект стремится соответствовать современным стандартам разработки программного обеспечения.

<details>
    <summary><b>Показать</b></summary>

### 🏗️ Архитектурные принципы:
- **Чистая архитектура**: Строгое разделение между слоями данных, бизнес-логики и представления обеспечивает масштабируемость и простоту поддержки
- **Dependency Injection**: Использование внедрения зависимостей делает код модульным и легко тестируемым
- **Repository Pattern**: Абстрактный интерфейс хранилища позволяет легко заменять источники данных (PostgreSQL, in-memory)
- **Publish-Subscribe Pattern**: Применение паттерна Pub/Sub для асинхронной передачи данных между компонентами через централизованного брокера, обеспечивая полную декомпозицию отправителей и получателей

### 💼 Бизнес-ориентированный подход:
- **Готовность к высоким нагрузкам**: Архитектура и выбранные технологии обеспечивают хорошую производительность при масштабировании
- **Поддержка микросервисной архитектуры**: Сервис легко интегрируется в микросервисную экосистему
- **Расширяемость**: Модульная структура позволяет легко добавлять новые функции и интегрироваться с другими системами

### 🛠️ Методологии разработки:
- **Code Generation**: Автоматическая генерация кода с помощью gqlgen минимизирует ручное написание повторяющегося кода
- **Database Migrations**: Структурированные миграции базы данных обеспечивают надежное обновление схемы
- **Environment Configuration**: Гибкая настройка через переменные окружения для различных сред развертывания
- **Детальное логирование**: Структурированные логи для мониторинга и диагностики системы

</details>

## 🔧 Стек:
  * Go 1.24
  * GraphQL (gqlgen)
  * PostgreSQL 17
  * Redis 9
  * Docker & Docker Compose
  * Gorilla WebSockets

## 🧪 Тестирование
Для проверки покрытия кода тестами:
```bash
go test ./internal/handlers/... -coverprofile cover.out 
```

```bash
go tool cover -func cover.out
```

![Покрытие бизнес логики юнит тестами](image.png)
(на снимке показан процент покрытия бизнес логики юнит тестами)

## 📜 Лицензия

Данный проект распространяется под лицензией Apache License 2.0. Это свободная лицензия с открытым исходным кодом, которая позволяет использовать, модифицировать и распространять код как в коммерческих, так и в некоммерческих целях.

Полный текст лицензии доступен в файле [LICENSE](./LICENSE).

```
Copyright 2023 Ozon Habr API

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```