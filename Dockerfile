FROM golang:onbuild
RUN go get bitbucket.org/liamstask/goose/cmd/goose
RUN apt-get update && apt-get install -y vim
RUN apt-get install -y cron
EXPOSE 4000