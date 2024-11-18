CREATE TABLE IF NOT EXISTS user_completed_tasks (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE
);