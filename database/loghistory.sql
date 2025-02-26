CREATE TABLE log_histories (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    command TEXT NOT NULL,
    count INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);