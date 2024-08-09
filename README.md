# Backend сервис для добавления новостей, написанный на Go
## Развертывание
### Общие настройки
Общие настройки приложения содержатся в [конфиге](https://github.com/dharmata314/news-service/blob/main/config/config.yaml). В зависимости от способа развертывания какие-либо параметры могут меняться. 
В конфиге содержатся основные данные, необходимые для работы приложения.

### Docker 
Для развертывания в Docker Compose создан файл [docker-compose.yaml](https://github.com/dharmata314/news-service/blob/main/docker-compose.yaml)
Необходимо запустить команду
```
docker compose up --build app
```
### Нативно
Для нативного запуска достаточно запустить приложение из папки [cmd](https://github.com/dharmata314/news-service/tree/main/cmd). 
Предварительно, необходимо установить зависимости из [go.mod](https://github.com/dharmata314/news-service/blob/main/go.mod) и изменить в [конфиге](https://github.com/dharmata314/news-service/blob/main/config/config.yaml) ```host: postgres``` на ```host: localhost```
## Общее
Приложение представляет из себя добавления новостей и получения списка новостей. 

Сервис устроен следующим образом:

Пользователи региструются в сервисе, авторизуются, затем добавляются новости, после чего можно получать список всех новостей или изменять существующие новости.

Доступны следующие операции:
 - Регистрация пользователя
 - Авторизация пользователя
 - Изменение данных пользователя
 - Добавление новости
 - Получение списка всех новостей
 - Изменение новости

Для доступа к большинству функционала (кроме регистрации и авторизации) необходим доступ по токену.
Токен выдается пользователю после авторизации.
В дальнейшем токен должен передаваться вместе с заголовком запроса:
```
Authorization: Bearer <token>
```

## Примеры запросов

Регистрация пользователя:
```
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://localhost:8080/users/new
```
Авторизация пользователя:
```
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://localhost:8080/login
```
Изменение данных пользователя:
```
curl -X PATCH \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"email": "newEmail@email.com", "password": "NewPassword", "id": 1}' \
http://localhost:8080/users/{id}
```
Добавление новости:
```
curl -X POST \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"Title": "Name", "Content": "Some content", "Categories": [1,2,3]}' \
http://localhost:8080/news
```
Получение списка всех новостей
```
curl -X GET \
-H "Authorization: Bearer <token>" \
http://localhost:8080/list
```
Изменение новости:
```
docker-compose exec app curl -X PATCH \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"Id": 1, "Title": "New_Name", "Content": "New_Content", "New_Categories": [1,2,3]}' \
http://localhost:8080/news/edit/{id}
```
