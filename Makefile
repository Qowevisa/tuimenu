def: test

test:
	go build -o ./bin/$@ ./examples/$@
