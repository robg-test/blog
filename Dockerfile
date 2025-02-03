FROM golang:1.23

WORKDIR /app

COPY . /app

RUN go install github.com/a-h/templ/cmd/templ@latest && templ generate && go build -o blog

ENV ENV=production

EXPOSE 443

CMD ["./blog"]
