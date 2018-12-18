protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  --go_out=plugins=grpc:. \
  --swagger_out=logtostderr=true,allow_delete_body=true:../swagger \
  --grpc-gateway_out=logtostderr=true,allow_delete_body=true:. \
  ./*.proto

# Need to figure out relative paths.
# Need to figure out how to rename the swagger to name-dev.

# YOU NEED TO ADD THESE INTO THE GENERATED SWAGGER IN security/firebase
# "x-google-issuer":"https://securetoken.google.com/phdigidev",
# "x-google-jwks_uri":"https://www.googleapis.com/service_accounts/v1/metadata/x509/securetoken@system.gserviceaccount.com",
# "x-google-audiences":"phdigidev"