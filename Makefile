.PHONY: proto
proto:
	@cd proto
	@protoc --go_out=plugins=grpc:. ./proto/*.proto