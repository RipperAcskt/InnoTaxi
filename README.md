# InnoTaxiUser API 

A user microservice for InnoTaxi app which provide possibility for working with user.

`app.env` is a configuration for service.

## Run the app

    go build -o ./bin/main ./cmd/main.go
    ./bin/main

## Run the tests

    go test ./internal/service 

# REST API

The REST API for service.

## Get profile by provided user's id

### Request

`GET /users/profile/:id`

    curl -i -H 'Accept: application/json' -H 'Authorization: Bearer ${access_token}' http://localhost:8080/users/profile/1

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Date: Fri, 17 Feb 2023 08:04:01 GMT
    Content-Length: 55

    {"name":"2","phone_number":"2","email":"2","raiting":0}

## Update profile by provided user's id

### Request

`PUT /users/profile/:id`

    curl -i -H 'Accept: application/json' -H 'Authorization: Bearer ${access_token}' -H 'Content-Type: application/json' -d "{\"phone_number\": \"+77777778\",\"email\": \"ripper@mail.ru\"}" -X PUT http://localhost:8080/users/profile/1

### Response

    HTTP/1.1 200 OK
    Date: Fri, 17 Feb 2023 08:16:11 GMT
    Content-Length: 0

## Selete profile by provided user's id

### Request

`DELETE /users/:id`

    curl -i -H 'Accept: application/json' -H 'Authorization: Bearer ${access_token}' -XDELETE http://localhost:8080/users/1

### Response

    HTTP/1.1 200 OK
    Date: Fri, 17 Feb 2023 08:20:07 GMT
    Content-Length: 0

## Create account

### Request

`POST /users/auth/sing-up`

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -d "{\"name\": \"2\",\"phone_number\": \"2\",\"email\": \"2\",\"password\": \"2\"}" -X POST http://localhost:8080/users/auth/sing-up

### Response

    HTTP/1.1 201 Created
    Date: Fri, 17 Feb 2023 08:26:40 GMT
    Content-Length: 0

## Sing in into accout

### Request

`POST /users/auth/sing-in`

    curl -i -H 'Accept: application/json' -H 'Content-Type: application/json' -d "{\"phone_number\": \"2\",\"password\": \"2\"}" -X POST http://localhost:8080/users/auth/sing-in

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Set-Cookie: refresh_token=${refresh_token}; Path=/users/auth; Max-Age=2592000; HttpOnly
    Date: Fri, 17 Feb 2023 08:29:26 GMT
    Content-Length: 159

    {"access_token":${access_token}}

## Update access token

### Request

`GET /users/auth/refresh`

    curl -i -H 'Accept: application/json' -H 'Authorization: Bearer ${access_token}' --cookie "refresh_token=${refresh_token}" -X GET http://localhost:8080/users/auth/refresh

### Response

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8
    Set-Cookie: refresh_token=${refresh_token}; Path=/users/auth; Max-Age=2592000; HttpOnly
    Date: Fri, 17 Feb 2023 08:33:56 GMT
    Content-Length: 159

    {"access_token":${access_token}}

## Logout from account

### Request

`GET /users/auth/logout`

    curl -i -H 'Accept: application/json' -H 'Authorization: Bearer ${access_token}' --cookie "refresh_token=${refresh_token}" -X GET http://localhost:8080/users/auth/logout

### Response

    HTTP/1.1 200 OK
    Set-Cookie: refresh_token=; Path=/users/auth; Max-Age=6; HttpOnly
    Date: Fri, 17 Feb 2023 08:37:06 GMT
    Content-Length: 0