-- Пользователи
CREATE TABLE users (
    id TEXT PRIMARY KEY, -- UUID
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Сессии (cookie-based)
CREATE TABLE sessions (
    id TEXT PRIMARY KEY, -- UUID
    user_id TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Посты
CREATE TABLE posts (
    id TEXT PRIMARY KEY, -- UUID
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Категории
CREATE TABLE categories (
    id TEXT PRIMARY KEY, -- UUID
    name TEXT NOT NULL UNIQUE
);

-- Привязка категорий к постам (многие ко многим)
CREATE TABLE post_categories (
    post_id TEXT,
    category_id TEXT,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Комментарии
CREATE TABLE comments (
    id TEXT PRIMARY KEY, -- UUID
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Лайки/дизлайки постов
CREATE TABLE post_reactions (
    user_id TEXT,
    post_id TEXT,
    reaction INTEGER CHECK(reaction IN (1, -1)),
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- Лайки/дизлайки комментариев
CREATE TABLE comment_reactions (
    user_id TEXT,
    comment_id TEXT,
    reaction INTEGER CHECK(reaction IN (1, -1)),
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);