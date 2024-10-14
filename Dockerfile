FROM golang:1.9

ENV PORT=80
ENV GIN_MODE=release
ENV WATCH_OTV=1
ENV DATABASE_DSN=tvthailand:A4xA9DN46cyGElgR@tcp(db.tvthailand.me:3306)/tvthailand?parseTime=true
ENV REDIS_HOST=db.tvthailand.me

WORKDIR /go/src/github.com/code-mobi/tvthailand.me
COPY . .

# RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run", "tvthailand.me"]

EXPOSE 3306
EXPOSE 6379
EXPOSE 80
