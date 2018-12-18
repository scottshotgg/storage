ENV=dev \
go build -o ./dist/storage \
&& \
./dist/storage both \
--servicename=storage \
--server-ip="localhost" \
--rest-port="6000" \
--rpc-addr="localhost" \
--rpc-port="6001"
