FROM golang:1.24

WORKDIR /app

COPY . /app

RUN go get -tool github.com/a-h/templ/cmd/templ@latest && go tool templ generate && go build -o blog

ENV ENV=production

EXPOSE 443

CMD ["./blog"]
