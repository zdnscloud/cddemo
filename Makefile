GOSRC = $(shell find . -type f -name '*.go')

build: cddemo 

cddemo: $(GOSRC) cddemo.go
		go build cddemo.go

test:
		go test -v -timeout 60s -race ./...

docker:
		docker build -t zdnscloud/cddemo:v0.1.0 .
			docker image prune -f
				docker push zdnscloud/cddemo:v0.1.0

clean:
		rm -rf cddemo 

.PHONY: clean install
