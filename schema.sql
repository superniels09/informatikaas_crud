CREATE TABLE IF NOT EXISTS users (
    id           BIGINT      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username     VARCHAR(255) NOT NULL UNIQUE,
    passhash     VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255),
    bio          TEXT,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS games (
    id           BIGINT      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    description  TEXT NOT NULL,
    script       BLOB NOT NULL,
    user_id      BIGINT      NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX (user_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id         BIGINT      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    content    TEXT NOT NULL,
    user_id    BIGINT      NOT NULL,
    game_id    BIGINT      NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (game_id) REFERENCES games(id),
    INDEX (game_id, created_at)
);

CREATE TABLE IF NOT EXISTS favorites (
    user_id    BIGINT      NOT NULL,
    game_id    BIGINT      NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, game_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (game_id) REFERENCES games(id),
    INDEX (game_id)
);
