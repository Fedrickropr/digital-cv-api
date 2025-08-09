To start:

DB:
```
docker compose up
```

Dev build (with hot reload)
```
CompileDaemon 
```
CompileDaemon --build="go build -o digital-cv-api.exe" --command="./digital-cv-api.exe" --directory="." --recursive
```