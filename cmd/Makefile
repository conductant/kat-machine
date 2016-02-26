all: clean kat-machine

clean:
	-rm kat-machine

# Common makefile that extracts git version info and generates the LDFLAGS variable.
include ../hack/make/version.mk

KEY_NAME?=example_key

gen-key-pair:
	openssl genrsa -out $(KEY_NAME) 4096
	openssl rsa -pubout -in $(KEY_NAME) -out "$(KEY_NAME).pub"

kat-machine:
	${GODEP} go build -ldflags "$(LDFLAGS)" kat-machine.go

kat-machine-run: kat-machine
	./kat-machine --logtostderr run \
	--port=7070 \
	--public_key_url=file://./example_key.pub

kat-machine-quit: kat-machine
	curl -i -X POST \
	-H "Authorization: Bearer `./kat-machine token -private_key_url=file://./example_key -scope=admin-kill`" \
	http://localhost:7070/quitquitquit
