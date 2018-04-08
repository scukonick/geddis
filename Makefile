

build:
	go build -o geddis main.go

run-local:
	go run main.go

docker:
	go build -o geddis main.go
	docker build -t scukonick/geddis .
