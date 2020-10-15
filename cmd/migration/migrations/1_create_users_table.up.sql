CREATE TABLE IF NOT EXISTS tasks (
                       id serial  PRIMARY KEY,
                       exec_time    timestamp,
                       name varchar(120),
                       created_at timestamp
);

