FROM golang

ENV GIN_MODE=release
ENV WATCH_OTV=1
ENV DATABASE_DSN=makathon:G00UltraMrds@tcp(makathoninstance.c2ckzrktsntv.us-east-1.rds.amazonaws.com:3306)/tvthailanddb?parseTime=true
ENV REDIS_HOST=tvthailand.gntesa.0001.use1.cache.amazonaws.com

ADD . /go/src/github.com/code-mobi/tvthailand.me
WORKDIR /go/src/github.com/code-mobi/tvthailand.me
RUN go build

ENTRYPOINT /go/src/github.com/code-mobi/tvthailand.me/tvthailand.me

EXPOSE 3000
