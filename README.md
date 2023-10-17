# Company API project


## How to run program

#### Using docker compose
```sh
cd docker
docker compose build
docker compose up
```

#### Run localy
```sh
DATABASE_URL=root:@tcp(127.0.0.1:3306)/xmproject go run internal/cmd/main.go 
```

## Database migration
Project uses mysql database. Migration files are available in directory `migrations`.
You an use https://github.com/golang-migrate/migrate to run migration.

## OpenAPI 
OpenAPI can be found in `openapi/xmproject.yml`
Link to swagger editor https://app.swaggerhub.com/apis/JANTOMX/CompanyAPI/1.0.0

## JWT Authentication
The API validates a JWT token received in an HTTP request using the `HS256` algorithm. The default private key is set to `secret`, but you can customize it using the `JWT_PRIVATE_KEY` environment variable. Please note that additional authorization using custom JWT claims has not been implemented at this time.

Example request for `secret` private key:
```
curl --location 'http://127.0.0.1:8080/companies' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.e30.XmNK3GpH3Ys_7wsYBfq4C3M6goz71I7dTgUkuIa5lyQ' \
--data '{
    "name": "some_name",
    "employees_count": 7,
    "registered": true,
    "type": "cooperative"
}'
```

