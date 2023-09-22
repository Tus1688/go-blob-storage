FROM golang:alpine AS builder

RUN apk update && apk add alpine-sdk git gcc && rm -rf /var/cache/apk/*
RUN mkdir -p api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app ./main.go

FROM nginx:stable-alpine3.17-slim
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /api/app .
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d
COPY script.sh /script.sh
RUN chmod +x /script.sh

EXPOSE 5000

CMD ["/script.sh"]