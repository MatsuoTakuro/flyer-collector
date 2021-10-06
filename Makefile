run:
	go build && go run flyer-collector

st2:
	go build && go run flyer-collector -store=2

rmf:
	rm -rf files/*
