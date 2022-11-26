mod:
	go mod download

build-eavesdropper:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o eavesdropper .

test:
	go test ./...

lint:
	golint ./...

check-suite: test
