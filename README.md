# Una api sencilla de practica
Es una api sencilla que es escalable y sigue buenas practicas , ademas de contar con la documentacion de swagger esta api y la arquitectura se basa en este video 
https://youtu.be/OGhQhFKvMiM?t=9475

## El archivo de configuracion seria un yaml este puede verse algo asi
y bueno recomiendo que este en una carpeta llamada config algo asi seria la estructura

cmd
--students-api
----main.go
config
--local.yaml
internal

```yaml
env: "dev"
storage_path: "storage/storage.db" #o tu conexion si es que no tienes pensado usar sqlite
http_server:
  address: "localhost:8082"
smtp:
  from: "youremail@email.com"
  password: "yourpassword"
  smtpHost: "yourhost"
  smtpPort: "yourport"
  
```