.PHONY: run-app

dev:
	npx tailwind -i ./web/static/css/input.css -o ./web/static/css/output.css --watch &
	air

run:
	go build -o app .
	./app
