protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  --go_out=plugins=grpc:. \
  --swagger_out=logtostderr=true,allow_delete_body=true:../swagger \
  --grpc-gateway_out=logtostderr=true,allow_delete_body=true:. \
  ./*.proto \
&& \
# Generate the easyjson fast un/marshal bindings
easyjson -all *.pb.go
# Need to figure out relative paths.
# Need to figure out how to rename the swagger to name-dev.