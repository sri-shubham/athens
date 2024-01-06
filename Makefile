build:
	go build -o athens main.go

clean:
	rm -f athens

run: build
	docker-compose up -d
	./athens

stop:
	docker-compose down
