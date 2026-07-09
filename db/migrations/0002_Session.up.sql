CREATE TABLE IF NOT EXISTS session (session_key SERIAL PRIMARY KEY,user_id INT NOT NULL,expire_at INT);
