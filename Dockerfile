# TODO
FROM golang

WORKDIR /app

ENV DBUSER=root
ENV DBPASS=superu
ENV GOPATH=

COPY . .

WORKDIR /app/main

EXPOSE 80

CMD [ "go", "run", "/app/main/server.go" ]
