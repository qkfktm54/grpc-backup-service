protoc --go_out=./pkg/backup      \
 --go_opt=paths=source_relative      \
 --go-grpc_out=./pkg/backup \
 --go-grpc_opt=paths=source_relative \
 ./proto/*.proto