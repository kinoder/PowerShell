CREATE TABLE log_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    command VARCHAR(255),
    count INT DEFAULT 1
);