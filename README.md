# Schedule

# For SocialTech

1) run docker-compose.yaml
2) copy .env.dist, paste, and change file name .env
3) create database
3) run cmd/migration/main.go
4) run cmd/scheduler/main.go
5) run pkg/script.go
#
 * localhost:SERV_ADDR/tasks - add task, body pkg/scheduler/task.go
 * localhost:SERV_ADDR/tasks/{id} - remove task
 
#
> project layout: https://github.com/golang-standards/project-layout
>
* реализация через полинг базы данных и очередь 
* скалирование можно сделать создав несколько инстансов worker_poll.go с одним manger.go для полинга

#Будущее:
1) можно сделать скалирование на уровне сервиса, улучшив работу с базой
2) улучшить amqplb
3) добавить интеграционные тесты
4) юнит тесты для воркер пула
