build:
	docker build -t goloadtester .

run:
	docker run goloadtester --url=http://example.com --requests=1000 --concurrency=20

test:
	go test -v





