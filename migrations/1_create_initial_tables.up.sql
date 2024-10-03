CREATE TABLE IF NOT EXISTS templates (
    id BIGSERIAL PRIMARY KEY,
    template_name VARCHAR(255) UNIQUE NOT NULL,
    template_content TEXT NOT NULL,
    required_parameters VARCHAR(255)[]
);

CREATE TABLE IF NOT EXISTS notification_history (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255),
    text TEXT,
    status VARCHAR(255),
    error_message TEXT
);