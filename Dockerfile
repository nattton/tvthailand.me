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

# docker-machine create --driver digitalocean --digitalocean-access-token=1cfcb82e4e2551e01a90afd93d2f176e70a904e7bddfd0c22dab40095e2dd75e --digitalocean-region "sgp1" tvthailand-me
# docker build -t tvthailand-me .
# docker run -it -d --publish 80:80 --name tvthailand_web  --restart always tvthailand-me
# docker stop tvthailand_web;docker run -it -d --publish 80:80 --name tvthailand_web  --restart always tvthailand-me