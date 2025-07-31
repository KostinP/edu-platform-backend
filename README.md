# edu-platform-backend

Содержание папки pkg:
- pkg/telegram  
  - bot.go  
  - oauth.go  
- pkg/middleware  
  - context_keys.go  
  - jwt.go  
  - link_visitor_user.go  
  - set_user_id.go  
  - visitor.go  
- pkg/logger/logger.go  
- pkg/geo/geo.go  
- pkg/entity/base.go  
- pkg/db  
  - redis.go  
  - postgres.go  
- pkg/config  
  - types.go  
  - config.go  

---

## pkg/telegram

### bot.go

Файл отвечает за инициализацию и создание нового клиента Telegram Bot API.

Функция:

```go
func New(token string) (*tgbotapi.BotAPI, error)
```

* Принимает токен бота.
* Создаёт новый экземпляр клиента Telegram.
* Включает режим отладки (debug = true).
* Логирует успешную авторизацию с именем бота.

---

### oauth.go

Реализация Telegram OAuth авторизации:

* `AuthData` — структура для данных, полученных из Telegram OAuth (id пользователя, username, имя, hash).
* `ParseTelegramAuth(r *http.Request) AuthData` — парсит и извлекает параметры авторизации из HTTP-запроса.
* `VerifyTelegramAuth(authData AuthData, botToken string) bool` — проверяет корректность данных с помощью HMAC SHA-256, используя секрет бота.

---

## pkg/middleware

### pkg/midddleware/context\_keys.go

* Определяет ключи для хранения данных в контексте Echo.

```go
const UserIDKey = "user_id"
```

---

### pkg/midddleware/jwt.go

Middleware для проверки JWT токена и управления сессиями пользователя.

* Проверяет заголовок Authorization на наличие JWT в формате Bearer.
* Валидирует подпись JWT с помощью секретного ключа.
* Извлекает из токена `user_id`, `session_id` и время истечения `exp`.
* Проверяет валидность и срок действия токена и сессии в базе.
* Обновляет время последней активности сессии.
* Если токен истекает менее чем через 24 часа, генерирует новый и возвращает его в заголовке `X-Refresh-Token`.
* Устанавливает `user_id` и `session_id` в контекст Echo для дальнейшего использования.

---

### pkg/midddleware/link\_visitor\_user.go

Middleware связывает гостевого посетителя (visitor\_id) с авторизованным пользователем (user\_id):

* Если есть оба идентификатора, вызывает метод `LinkVisitorToUser` сервиса пользователя.
* Ошибки связывания не блокируют выполнение запроса — логируются при необходимости.

---

### pkg/midddleware/set\_user\_id.go

Middleware для установки `user_id` из заголовка HTTP `X-User-ID`:

* Если заголовок присутствует и содержит валидный UUID, устанавливает его в контекст.
* Иначе пропускает запрос без ошибок.

---

### pkg/midddleware/visitor.go

Middleware для работы с гостевыми посетителями:

* Проверяет наличие cookie с `visitor_id`.
* Если отсутствует — генерирует новый UUID, устанавливает cookie на 1 год, с флагами HttpOnly и SameSite=Lax.
* Помещает `visitor_id` в контекст Echo.

---

## pkg/logger/logger.go

Простой логгер с тремя уровнями:

* Info(msg string) — информационные сообщения в stdout с префиксом "INFO".
* Error(msg string, err error) — ошибки в stderr с префиксом "ОШИБКА".
* Fatal(msg string, err error) — критические ошибки, вызывающие завершение программы.

Использует стандартный пакет `log` с флагами даты, времени и имени файла.

---

## pkg/geo/geo.go

Заглушка геолокации по IP для локальной разработки.

* Функция `Lookup(ip string) (country, city string)` всегда возвращает `"Unknown", "Unknown"`.

---

## pkg/entity/base.go

Базовая структура сущности для моделей с общими полями:

```go
type Base struct {
    ID        uuid.UUID `json:"id" db:"id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
    AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
}
```

Методы:

* `Init(authorID uuid.UUID)` — инициализирует ID, CreatedAt, UpdatedAt и AuthorID.
* `Touch()` — обновляет поле UpdatedAt текущим временем.

---

## pkg/db

### redis.go

* Инициализация Redis клиента по URL.
* Функция `NewRedis(url string) (*Redis, error)` парсит URL и создаёт `redis.Client`.
* Структура `Redis` содержит поле `Client *redis.Client`.

---

### postgres.go

* Подключение к PostgreSQL через pgxpool.
* Функция `ConnectPostgres(cfg *config.Config) *pgxpool.Pool`:

  * Формирует DSN из конфигурации.
  * Пингует базу для проверки доступности.
  * При ошибках логирует и завершает работу.
  * Логирует успешное подключение.

---

## pkg/config

### types.go

Определяет алиасы для конфигурационных параметров:

* `BotToken string`
* `JwtSecret string`

---

### config.go

Структура `Config` описывает всю конфигурацию приложения:

* Поддерживает загрузку из YAML-файлов (`configs/dev.yaml`, `configs/prod.yaml` и т.д.) в зависимости от переменной окружения `APP_ENV` (по умолчанию `dev`).
* Включает конфигурацию серверной части, базы данных, Telegram и JWT.
* При ошибках загрузки конфигурации логирует и завершает приложение.

---

## pkg/pagination

### Общая информация

В проекте реализована пагинация, разделённая на два слоя:

* HTTP слой (`pkg/http/pagination`) — парсит параметры пагинации из HTTP-запроса.
* Доменный слой (`pkg/pagination`) — бизнес-логика и генерация SQL-запросов с пагинацией.

---

### HTTP слой — pkg/http/pagination

* Структура `PaginationQueryParams` содержит поля: Limit, Offset, SortBy, Order.
* Функция `ParsePaginationParams(c echo.Context) PaginationQueryParams` — извлекает параметры из query string и нормализует их.
* Метод `ToDomainParams()` конвертирует HTTP структуру в доменную `pagination.Params`.

---

### Доменный слой — pkg/pagination

* Структура `Params` хранит параметры пагинации.
* Метод `Normalize()` корректирует параметры, выставляя значения по умолчанию, если они невалидны.
* Функция `SQLWithPagination(baseQuery string, p Params, allowedSortFields map[string]string)` формирует SQL-запрос с LIMIT, OFFSET и ORDER BY.

---

### Пример использования пагинации

```go
func (h *TagHandler) ListTags(c echo.Context) error {
	// 1. Парсим HTTP параметры пагинации
	httpPag := myHttpPagination.ParsePaginationParams(c)

	// 2. Конвертируем HTTP DTO в доменный тип
	domainPag := httpPag.ToDomainParams()

	// 3. Используем доменные параметры в usecase для получения данных
	tags, total, err := h.TagUsecase.ListTags(c.Request().Context(), domainPag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list tags"})
	}

	// 4. Возвращаем ответ с пагинацией
	return c.JSON(http.StatusOK, dto.PaginatedResponse[*entity.Tag]{
		Items:  tags,
		Total:  total,
		Limit:  domainPag.Limit,
		Offset: domainPag.Offset,
	})
}
```

---

### Рекомендации по пагинации

* Используйте HTTP слой (`pkg/http/pagination`) только в хендлерах.
* В бизнес-логике и репозиториях работайте с `pkg/pagination.Params`.
* Не создавайте циклических импортов между слоями.
* Нормализуйте параметры перед использованием.
* Используйте `SQLWithPagination` для генерации безопасных SQL-запросов.
