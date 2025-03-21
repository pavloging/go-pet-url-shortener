# REST API сервис на Go

Что включает в себя этот проект?
- Взаимодейстие и разработка REST API
- Обращение к бд - SQL
- Router - go-chi/chi
- Логирование - slog
- Тесты - Unit (тесты обработчиков), Функциональные тесты
- CI/CD - Workflow GitHub (CI) и отправка на сервер (CD)

# С Начала до 21:34

## Структура проекта:
- cmd папка - внутри cmd лежат команды для взаимодействия с приложением (запуск и т.д.)
- config - папка по настройки параметров для нашего приложения
- internal - в этой папке располагается код, который НЕ должен быть доступениз других модулей (похоже на приватное хранилище)

## Важный момент!
**В реальных репозиториях файл `.env` НЕЛЬЗЯ пушить в открытый доступ на GitHub.** Поэтому я создал в корне проекта `example.env` в целях деменстрации полей, которые записаны в `.env`. Файл `example.env` будет в открытом доступе, поэтому не заполняйте его паролями и приватной информацией, он лишь нужен для деменстрации ключей которые есть в приложении, чтобы следующий программист мог вставить свои приватные данные (пароль от базы например) и так же пользоваться приложением.

- Что вам делать?
Вам нужно создать в корне проекта файл с названием: `.gitignore` и в нем написать - `.env`

- Зачем это делать?
Пути которые прописаны в `.gitignore` не будут видны Git, это сделано в целях безопастности, чтобы файлы с паролями не отправлять в открытый доступ

## Небольшой конспект видео:
- **Хорошая практика:** Приставка Must в функциях обычно используется, когда вместо возврата ошибки будет вызывать панику (Автор делает так при инициализации конфига)

- **Хорошая практика:** Удалять fmt с конфигом, так как в логи может залететь пароль в открытом виде и злоумышленик может ними воспользоваться

## Возможные ошибки:

### Если не работает `go get github.com/ilyakaznacheev/cleanenv`
Сначала инициализаируйте go.mod командной `go mod init url-shortener` и потом устанавливайте библиотеку

### Ошибка: CONFIG_PATH is not set
Автор хочет чтобы мы указывали путь до файла в теринале при каждом запуске. Но мы воспользуемся другим способом:
1. Создайте файл в корне проекта .env
2. Пришите в этот файл: CONFIG_PATH=./config/local.yaml
3. Добавьте библиотку для удобной работы с .env: go get github.com/joho/godotenv 
4. Примените библиотеку:
```go
func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

    // ...
```

### Ошибка: cannot read config: field
Переходим в config -> local.yaml и добавляем user и password
Например файл получиться:
```yaml
env: "local" # local, dev, prod
storage_path: "./storage/storage.db"
http_server:
  address: "localhost:8082"
  timeout: 4s
  idle_timeout: 60s
  user: 123
  password: 123
```

# От 21:34 до 43:56

## Возможные ошибки:

### Ошибка: msg="failed to init storage" error="storage.sqlite.New: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub"

Вот <a href="https://github.com/mattn/go-sqlite3/issues/855#issuecomment-2267489894" target="_blank">ответ</a> на эту ошибку и на последующую. Но если вкртце, для Windows пропишите в терминале `go env -w CGO_ENABLED=1` и затем скачайте tdm64-gcc <a href="https://jmeubank.github.io/tdm-gcc/" target="_blank">последняя версия</a>, после чего перезапустите VSCode.
Для других OC решение описано по ссылке выше.

### Открываем файл .db
Вы можете самостоятельно скачать программу DBeaver или скачать расшерение SQLite <a href="https://www.youtube.com/watch?v=By-UUTO09xA" target="_blank">подробне</a>

# Метод DeleteURL
Сделал метод DeleteURL, который возвращает количество удалимых элементов после его работы. По факту, данное значение будет или 0 или 1. Так как элемента может не быть или он только один с таким названием, потому что названия в базе уникальные.

# Middleware

Без Middeware - Запрос -> Ответ
C Middeware - Запрос -> Middeware -> Ответ

То есть Middeware это промежуточный этап между запросом и ответом!

# Logger for Handler with Colors

На реальных проектах, рекомендуется использовать стандартный logger, а не реализацию со своим.