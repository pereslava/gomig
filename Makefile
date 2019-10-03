doc:
	docker-compose -f godoc-compose.yml up --build -d

stopdoc:
	docker-compose -f godoc-compose.yml down -v

test:
	go test