$file = $args[0] + ".proto"
Write-Output "Generating $file"
protoc --go_out=../internal/api/types --go_opt=paths=source_relative --go-grpc_out=../internal/api/types --go-grpc_opt=paths=source_relative --proto_path=../internal/api/proto $file
