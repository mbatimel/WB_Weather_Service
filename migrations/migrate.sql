-- Создание таблицы для хранения информации о городах
CREATE TABLE cities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL
);

-- Создание таблицы для хранения предсказаний погоды
CREATE TABLE weather_forecasts (
    id SERIAL PRIMARY KEY,
    city_id INT NOT NULL REFERENCES cities(id) ON DELETE CASCADE,
    temp DECIMAL(5, 2) NOT NULL,
    date DATE NOT NULL,
    weather_data JSON NOT NULL,
    UNIQUE (city_id, date)
);

