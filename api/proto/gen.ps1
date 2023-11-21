$file = $args[0] + ".proto"
echo "Generating $file"
protoc --go_out=../types --go_opt=paths=source_relative --go-grpc_out=../types --go-grpc_opt=paths=source_relative .\$file
