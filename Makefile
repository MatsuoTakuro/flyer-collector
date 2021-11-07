run:
	go build && go run flyer-collector

test:
	go build && go run flyer-collector -store=2

rmf:
	rm -rf files/*
