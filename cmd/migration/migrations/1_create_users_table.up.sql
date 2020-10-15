CREATE TABLE IF NOT EXISTS tasks (
                       id serial  PRIMARY KEY,
                       exec_time    timestamp,
                       name varchar(120),
                       is_finished boolean default false,
                       in_progress boolean default false,
                       created_at timestamp
);

