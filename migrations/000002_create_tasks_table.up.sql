CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    uuid UUID,
    title VARCHAR(50),
    description TEXT,
    reward INT
);