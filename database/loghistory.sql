CREATE TABLE log_history (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    command TEXT NOT NULL,
    count INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);