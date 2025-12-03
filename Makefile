.PHONY: build run test clean

build:
	go build -o andi-custodian ./cmd/andi-custodian

run: build
	./andi-custodian

test:
	go test -v ./...

clean:
	rm -f andi-custodian