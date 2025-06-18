# News_site

## Описание проекта

News_site - это современный новостной портал, построенный на Go с использованием чистой архитектуры. Проект
предоставляет функционал для просмотра, создания и управления новостными статьями с системой аутентификации
пользователей.

## Основные возможности

- 🔐 Аутентификация пользователей (регистрация, вход, выход)
- 📰 Просмотр списка новостных статей
- ✍️ Создание новых статей
- 📖 Просмотр полного текста статей
- 🎨 Современный адаптивный интерфейс
- 🔒 Безопасное хранение паролей
- 🛡️ JWT-аутентификация
- 👍 Возможность ставить лайки к статьям

## Технологии

- **Backend:**
    - Go 1.24
    - MySQL
    - Gorilla Mux (маршрутизация)
    - JWT (аутентификация)
    - Bcrypt (хеширование паролей)

- **Frontend:**
    - HTML5
    - Bootstrap 5.3
    - CSS3
    - JavaScript

## Структура проекта

```
├── cmd/                    # Точка входа в приложение
│   └── main.go            # Основной файл приложения
│
├── internal/              # Внутренний код приложения
│   ├── auth/             # Аутентификация и авторизация
│   │   ├── jwt/         # JWT токены
│   │   └── middleware/  # Middleware для аутентификации
│   ├── config/          # Конфигурация приложения
│   ├── core/            # Ядро приложения
│   ├── db/              # Работа с базой данных
│   ├── handlers/        # HTTP обработчики
│   ├── models/          # Модели данных
│   └── services/        # Бизнес-логика
│
├── web/                  # Веб-интерфейс
│   └── templates/       # HTML шаблоны
│       ├── createPost.html   # Шаблон создания статьи
│       ├── editPost.html     # Шаблон редактирования статьи
│       ├── showPost.html     # Шаблон просмотра статьи (теперь с лайками)
│       ├── userPosts.html    # Шаблон статей пользователя
│       ├── login.html        # Шаблон входа
│       ├── register.html     # Шаблон регистрации
│       ├── header.html       # Шапка сайта
│       └── footer.html       # Подвал сайта
│
├── .env                  # Конфигурация окружения
├── .gitignore           # Игнорируемые файлы
├── go.mod               # Зависимости Go
├── go.sum               # Контрольные суммы зависимостей
└── README.md            # Документация проекта
```

## Установка и запуск

1. **Предварительные требования:**
    - Go 1.24 или выше
    - MySQL 5.7 или выше
    - Git

2. **Клонирование репозитория:**
   ```bash
   git clone https://github.com/yourusername/News_site.git
   cd News_site
   ```

3. **Установка зависимостей:**
   ```bash
   go mod download
   ```

4. **Настройка базы данных:**
    - Создайте базу данных MySQL
    - Создайте таблицу users:
   ```sql
   CREATE TABLE users (
       id INT AUTO_INCREMENT PRIMARY KEY,
       email VARCHAR(255) UNIQUE NOT NULL,
       password VARCHAR(255) NOT NULL,
       username VARCHAR(255) UNIQUE NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```
    - Создайте таблицу articles:
   ```sql
   CREATE TABLE articles (
       id INT AUTO_INCREMENT PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       anons TEXT NOT NULL,
       full_text TEXT NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```
    - Создайте таблицу likes:
   ```sql
   CREATE TABLE likes (
       id INT AUTO_INCREMENT PRIMARY KEY,
       user_id INT NOT NULL,
       article_id INT NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
       FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
   );
   ```

5. **Настройка окружения:**
   Создайте файл `.env` в корневой директории:
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   SERVER_PORT=8080
   JWT_SECRET=your_secret_key
   ```

6. **Запуск приложения:**
   ```bash
   go run cmd/main.go
   ```

7. **Доступ к приложению:**
   Откройте браузер и перейдите по адресу: `http://localhost:8080`

## API Endpoints

- `GET /` - Главная страница со списком статей
- `GET /login` - Страница входа
- `POST /login` - Обработка входа
- `GET /register` - Страница регистрации
- `POST /register` - Обработка регистрации
- `GET /logout` - Выход из системы
- `GET /create` - Страница создания статьи
- `POST /save_article` - Сохранение новой статьи
- `GET /post/{id}` - Просмотр отдельной статьи

## Безопасность

- Все пароли хешируются с использованием bcrypt
- JWT токены для аутентификации
- Защита от SQL-инъекций
- Безопасное хранение конфигурации
- HTTP-only cookies для токенов

## Поддержка

`
Если у вас возникли вопросы или проблемы, создайте issue в репозитории проекта.`
