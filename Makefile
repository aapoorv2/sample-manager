
generate_grpc:
    # make add_to_path
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/sample-manager.proto

add_to_path:
    export PATH="$PATH:$(go env GOPATH)/bin"

start:
    go run server/server.go

test:
    go test ./... -v