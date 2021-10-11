gen: 
	protoc --go_out=./proto --proto_path=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative ./proto/upload.proto