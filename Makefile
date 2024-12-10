BINARY_NAME=gobatmon

build:
	go mod tidy
	go build -o ${BINARY_NAME} main.go

clean:
	go clean
	rm -f ${BINARY_NAME}
