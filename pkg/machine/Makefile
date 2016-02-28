all: test-machine

test-machine:
	${GODEP} go test ./...  -logtostderr -check.vv -v ${TEST_ARGS}
