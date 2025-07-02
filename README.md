#  Friends!💵🫂⏳...⌛🚶‍♂️ 🏃‍➡️
## Сервис по поиску друзей на сутки🙋
CRUD c аутентификацией и авторизацией.Реализован согласно принципам REST API.

## Установка
```shell
git clone git@github.com:AndreySirin/Friends.git 
```
## Перед запуском рекомендуется установить утилиты необходимые для проверки кода и запустить их командами из Makefile
```shell
make dev-tools
make lint
```
## Запуск docker compose
```shell
make up
```
## Удаление контейнеров docker compose
```shell
make down
```
# Описание методов API
### Приведенные примеры подразумевают отправку запросов с помощью Postman.
Методы CRUD доступны только аутентифицированным пользователям, чтоб воспользоваться ими 
в Postman во вкладке "Headers" в колонке "Key" указываем Authorization,
а в колонке "Value" — Bearer <токен>
```shell
пример:
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzUxNDY3NjQyLCJpYXQiOjE3NTE0NjY3NDJ9.v9H497yA8jCXHRKYdzE1m_V6W2Q55m_nxijFkIimbrI
```
## Регистрация пользователя
```shell
метод:POST
URL:http://localhost:8080/api/v1/registration
body:
{
    "name":"andrey",
    "email":"andrey@email.com",
    "password":"qwerty123"
}
ответ:"Registration successful"
status:200
```

## Аутентификация пользователя
```shell
метод:POST
URL:http://localhost:8080/api/v1/authentication
body:
{
    "email":"andrey@email.com",
    "password":"qwerty123"
}
ответ:
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzUxNDY3NjQyLCJpYXQiOjE3NTE0NjY3NDJ9.v9H497yA8jCXHRKYdzE1m_V6W2Q55m_nxijFkIimbrI"
}
status:200
```
## Обновление токена 
```shell
метод:GET
URL:http://localhost:8080/api/v1/refreshToken
body:
{}
ответ:
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzUxNDY3NjQyLCJpYXQiOjE3NTE0NjY3NDJ9.v9H497yA8jCXHRKYdzE1m_V6W2Q55m_nxijFkIimbrI"
}
status:200
```
## Запрос на наличие анкет
```shell
метод:GET
URL:http://localhost:8080/api/v1/price
body:
{}
ответ:html страница price.html
status:200
```
## Запрос на создание анкеты
```shell
метод:POST
URL:http://localhost:8080/api/v1/user
body:
{
 "name":"Leo",
 "hobby":"footbal",
 "price":100
}
ответ:
{
     "id": "6",
     "message": "Product added successfully"
}
status:200
```
## Запрос на обновление анкеты
```shell
метод:PUT
URL:http://localhost:8080/api/v1/user
body:
{
"id":6,
"name":"Leo Messi",
"hobby":"football",
"price":100
 }
ответ:
{
     "id": "6",
     "message": "Product updated successfully"
}
status:200
```
## Запрос на удаление анкеты
```shell
метод:DELETE
URL:http://localhost:8080/api/v1/user/{ID}
body:
{}
ответ:
{
    "message": "Product deleted successfully"
}
status:200
```



