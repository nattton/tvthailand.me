FROM golang

ENV DATABASE_DSN=makathon:G00UltraMrds@tcp(makathoninstance.c2ckzrktsntv.us-east-1.rds.amazonaws.com:3306)/tvthailanddb?parseTime=true

ADD . /go/src/github.com/code-mobi/tvthailand.me
WORKDIR /go/src/github.com/code-mobi/tvthailand.me
RUN go build

ENTRYPOINT /go/src/github.com/code-mobi/tvthailand.me/tvthailand.me

EXPOSE 3000
