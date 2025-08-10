# REST API for Digital CV


## Getting started
### Database
This project uses a docker Postgres database. If you have docker installed, you can use the included docker-compose. This starts a docker postgres image on port 5431.
DB:
```
docker compose up
```

### Run development
To run the development enviroment, use below command.
CompileDaemon
```
CompileDaemon --build="go build -o digital-cv-api.exe" --command="./digital-cv-api.exe" --directory="." --recursive
```