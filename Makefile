SRC=$(wildcard *.go)
PROJECT=unirestgo

.PHONY: test unirestgo

unirestgo: ${SRC}
	go build ${SRC}

test: ${SRC}
	go test ${SRC}

clean:
	${RM} ${PROJECT}