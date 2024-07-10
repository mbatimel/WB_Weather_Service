# Документация по проекту WB Weather Service

## Содержание
1. [Введение](#введение)
2. [Структура проекта](#структура-проекта)
3. [Зависимости](#зависимости)
4. [Сборка и запуск сервиса](#сборка-и-запуск-сервиса)
5. [API Эндпоинты](#api-эндпоинты)
6. [Примеры запросов](#примеры-запросов)

## Введение

WB Weather Service — это сервис для получения и обновления прогноза погоды для различных городов. Сервис взаимодействует с базой данных PostgreSQL и внешним API OpenWeather для получения данных о погоде.

## Структура проекта

Проект состоит из следующих основных частей:

- **Docker**: Используется для контейнеризации приложения и базы данных.
- **Go**: Основной язык программирования для разработки сервиса.
- **PostgreSQL**: База данных для хранения информации о городах и погодных прогнозах.
- **API**: Эндпоинты для взаимодействия с клиентами.

## Зависимости

Для запуска проекта необходимы следующие зависимости:

- Docker и Docker Compose
- Golang
- PostgreSQL

## Сборка и запуск сервиса

### Шаг 1: Клонирование репозитория

Сначала клонируйте репозиторий с проектом:

```bash
git clone https://github.com/mbatimel/WB_Weather_Service.git
cd WB_Weather_Service
```

### Шаг 2: Создание файла конфигурации

Создайте файл `config/config.yaml` с содержимым:

```yaml
server:
  host: "localhost"
  port: "8080"

repo:
    database: WB_developer
    user: mbatimel
    password: wb_il
    host: "localhost"
    port: "5432"
```

### Шаг 3.1: Запуск с использованием Docker(Данный шаг в разработке и в устранение демонов)

Запустите сервис с использованием Docker Compose:

```bash
docker-compose up --build
```

Docker Compose создаст и запустит два контейнера:
- `postgres`: Контейнер с базой данных PostgreSQL.
- `app`: Контейнер с приложением на Go.
### Шаг 3.2: Запуск сервера с использованием Makefile(работает 100%)

Запуск контейнера:
```bash
make up
```

Запуск сервера:
```bash
make servre
```
После завершения работы с сервером:
```bash
make down
```

### Шаг 4: Применение миграций

После запуска контейнеров примените миграции для базы данных:

```bash
docker-compose exec app ./migration
```

## API Эндпоинты

### Получение всех городов

**URL**: `/allCyties`

**Метод**: GET

**Описание**: Возвращает список всех городов.

**Пример ответа**:

```json
[
  "Belfast",
  "Berlin",
  "Cardiff",
  "Chelyabinsk",
  "Edinburgh",
  "Ekaterinburg",
  "Glasgow",
  "Kazan",
  "Liverpool",
  "London",
  "Manchester",
  "Moscow",
  "New York",
  "Novosibirsk",
  "Omsk","Oslo",
  "Paris",
  "Saint-Petersburg",
  "Samara",
  "Tokyo"
  ]
```

### Краткая информация о городе

**URL**: `/hortInfoCity?city={cityName}`

**Метод**: GET

**Параметры**:
- `city`: Название города.

**Описание**: Возвращает краткую информацию о городе, включая среднюю температуру и доступные даты прогноза.

**Пример запроса**:

```http
GET /hortInfoCity?city=London
```

**Пример ответа**:

```json
{"available_dates":[
  "2024-07-10",
  "2024-07-11",
  "2024-07-12",
  "2024-07-13",
  "2024-07-14"],
  "avg_temp":294.078,
  "city":"Berlin",
  "country":"DE"
}
```

### Полная информация о погоде в городе

**URL**: `/fullInfoCity?city={cityName}&date={date}`

**Метод**: GET

**Параметры**:
- `city`: Название города.
- `date`: Дата прогноза в формате `YYYY-MM-DD`.

**Описание**: Возвращает полную информацию о погоде в указанном городе на указанную дату.

**Пример запроса**:

```http
GET /fullInfoCity?city=London&date=2024-07-10
```

**Пример ответа**:

```json
 "city":"London",
  "country":"GB",
  "date":"2024-07-10T00:00:00Z",
  "temp":288.62,
  "weather_data":{
    "clouds":{"all":59},
    "dt":1720645200,
    "dt_txt":"2024-07-10 21:00:00",
    "main":{
      "feels_like":288.08,
      "grnd_level":1013,
      "humidity":71,
      "pressure":1017,
      "sea_level":1017,
      "temp":288.62,
      "temp_kf":0,
      "temp_max":288.62,
      "temp_min":288.62
      },
      "pop":0,
      "sys":{
        "pod":"n"
        },
      "visibility":10000,
      "weather":[{
        "description":"broken clouds",
        "icon":"04n",
        "id":803,
        "main":"Clouds"}],
      "wind":{"deg":259,"gust":7.78,"speed":3.61}}
```
### Добавление пользователя
**URL**: `/addNewUser?person={username}&password={pswd}`

**Метод**: ADD

**Параметры**:
- `person`: Инкнейм пользователя в системе.
- `password`: Пароль(хэшируется).

**Описание**: Добавляет пользователя в базу данных

**Пример запроса**:

```http
ADD /addNewUser?person=Mbatimel&password=1234
```

**Пример ответа**:

```json
"User created successfully"
```
### Добавление города в избранные у пользователя
**URL**: `/addCityToFuvorites?person={username}&password={pswd}&city={city}`

**Метод**: ADD

**Параметры**:
- `city`: Название города.
- `person`: Инкнейм пользователя в системе.
- `password`: Пароль(хэшируется).

**Описание**: Добавляет город в избранные у пользователя.

**Пример запроса**:

```http
ADD /addCityToFuvorites?person=Mbatimel&password=1234&city=Berlin
```

**Пример ответа**:

```json
"Added successfully"
```
### Полная информация об избранных городах
**URL**: `/favoritCityInfo?person={username}&password={pswd}`

**Метод**: GET

**Параметры**:
- `person`: Инкнейм пользователя в системе.
- `password`: Пароль(хэшируется).

**Описание**: Возвращает полную информацию по городам, которые находятся в избранных у пользователя

**Пример запроса**:

```http
GET /favoritCityInfo?person=Mbatimel&password=1234
```

**Пример ответа**:

```json
{
  "cities": [
        {
            "country": "RU",
            "latitude": 55.7504461,
            "longitude": 37.6174943,
            "name": "Moscow",
            "weather": {
                "date": "2024-07-10",
                "temp": 291.42,
                "weather_data": {
                    "clouds": {
                        "all": 2
                    },
                    "dt": 1720645200,
                    "dt_txt": "2024-07-10 21:00:00",
                    "main": {
                        "feels_like": 290.9,
                        "grnd_level": 1007,
                        "humidity": 61,
                        "pressure": 1025,
                        "sea_level": 1025,
                        "temp": 291.42,
                        "temp_kf": 0,
                        "temp_max": 291.42,
                        "temp_min": 291.42
                    },
                    "pop": 0,
                    "sys": {
                        "pod": "n"
                    },
                    "visibility": 10000,
                    "weather": [
                        {
                            "description": "clear sky",
                            "icon": "01n",
                            "id": 800,
                            "main": "Clear"
                        }
                    ],
                    "wind": {
                        "deg": 50,
                        "gust": 3.08,
                        "speed": 1.77
                    }
                }
            }
        }
  ]
}
```
# КОНЕЦ!!! А КТО ЧИТАЛ, МОЛОДЕЦ!!!!