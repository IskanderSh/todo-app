# REST API Для Создания TODO Списков на Go

### Для запуска приложения
```
docker pull postgres
docker run --name=todo-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres
make migrate
make run
```

### TODO
1) Запуск одной командой
2) Unit тесты
3) CI/CD