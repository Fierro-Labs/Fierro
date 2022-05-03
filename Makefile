build:
	cd ./src/ && go build -v -o .

daemon:
	ipfs daemon

run:
	go run ./src/

# ex: make select ARG1=TestAddFile
select: 
	go test -run $(ARG1) ./src/

test:
	go test -v ./src/
	
upload:
	go test -race -covermode=atomic -coverprofile=coverage.out ./src/
	bash <(curl -s https://codecov.io/bash) -t ${CODECOV_TOKEN}