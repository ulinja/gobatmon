BINARY_NAME=gobatmon

build:
	go mod tidy
	go build -o ${BINARY_NAME} gobatmon.go

clean:
	go clean
	rm -f ${BINARY_NAME}
