run:
	go run main.go model.go wordlist.go

wordterm:
	go build -o wordterm main.go model.go wordlist.go
	chmod +x wordterm

test:
	go test ./...

build:
	rm -rf dist
	mkdir dist
	go build -o dist/wordterm main.go model.go wordlist.go
