# Usamos la imagen de golang
FROM golang:1.24-alpine
# seleccionamos un directorio para nuestra app
WORKDIR /app

#copiamos los archivos .go
COPY go.* ./

#descargamos dependencias
RUN go mod download

#copiamos el resto de archivos
COPY . .
#definimos las variables de entorno para dev
ENV POSTGRES_URL = $POSTGRES_URL
#construimos la app
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
