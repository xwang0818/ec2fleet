export GOPATH=$(PWD)

all:
	go get github.com/aws/aws-sdk-go/service/ec2
	go build -o build/ec2fleet

test:
	go test src/util/*

clean:
	rm -rf build
