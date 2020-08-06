all: gen

.PHONY : gen

gen:
	protoc seaweed_raft.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative
