CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO users (email, hashed_password) 
VALUES 
    ('user1@example.com', '$2a$10$xJwL5vWZ5rY3VtG7pRwZ3.9z0tTk7Qe6z0tTk7Qe6z0tTk7Qe6z0tTk7'), -- Пароль: 123456
    ('user2@example.com', '$2a$10$yH8eL5vWZ5rY3VtG7pRwZ3.9z0tTk7Qe6z0tTk7Qe6z0tTk7Qe6z0tTk7'); -- Пароль: qwerty


INSERT INTO notes (user_id, title, body) VALUES
                                           (1, 'Покупки', 'Молоко, хлеб, яйца'),
                                            (1, 'Задачи на неделю', 'Закончить проект, позвонить клиенту'),
                                            (2, 'Идеи для блога', 'Написать статью о PostgreSQL, разобрать паттерны проектирования');
