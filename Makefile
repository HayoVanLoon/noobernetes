GO := /usr/local/go/bin/go
PROTOC := /usr/local/bin/protoc
DOCKER := /usr/bin/docker
MAKE := /usr/bin/make

GOOGLE_APIS := /home/hayo/Documents/googleapis
PROTO_SELF := v1

GO_GENERATED := ../go-generated
OUT := $(GO_GENERATED)/noobernetes/v1

TEST_ROOT := test
MOCK_TARGET := $(TEST_ROOT)/noobernetes/v1/noobernetes_mock.go


.PHONY: clean

clean:
	rm -rf $(OUT)

protoc: clean
	mkdir -p "$(OUT)"
	$(PROTOC) --go_out="plugins=grpc:$(OUT)" \
		--descriptor_set_out=api_descriptor.pb \
		-I"$(GOOGLE_APIS)" \
		-I"$(PROTO_SELF)" \
		v1/noobernetes.proto

build-counter: protoc
	$(DOCKER) build -t noobernetes_counter counter/.

build-ticker: protoc
	$(DOCKER) build -t noobernetes_ticker ticker/.

#
#mocks:
#	mkdir -p "$(TEST_ROOT)/$(ENVY_OUT)"
#	rm "$(MOCK_TARGET)"
#	mockgen -source envy/v1/envy.pb.go >> "$(MOCK_TARGET)"
#
#build:
#	$(go) build -i -o /dist/Envy_Server github.com/HayoVanLoon/envy/envy_server
#
#test-client: protoc
#	$(GO) build -i -o /tmp/___Envy_Client github.com/HayoVanLoon/envy/envy_client #gosetup
#	/tmp/___Envy_Client #gosetup
